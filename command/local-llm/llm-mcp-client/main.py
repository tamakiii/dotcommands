#!/usr/bin/env python3
"""
LLM MCP Client - Python implementation using LangChain

A client that enables local LLMs (via Ollama) to interact with MCP servers
using LangChain's native MCP integration for better maintainability.
"""

import asyncio
import sys
import argparse
from typing import List, Dict, Any, Optional, Union, Tuple
import subprocess
import json
import re

from langchain_ollama import ChatOllama
from langchain_core.messages import HumanMessage, SystemMessage, AIMessage, BaseMessage
from langchain_core.tools import BaseTool
from pydantic import BaseModel, Field, ConfigDict


class MCPToolArgs(BaseModel):
    """Base model for MCP tool arguments"""
    model_config = ConfigDict(extra='allow')


class MCPRequest(BaseModel):
    """Model for MCP JSON-RPC requests"""
    jsonrpc: str = "2.0"
    id: str
    method: str
    params: Dict[str, Any] = Field(default_factory=dict)


class MCPResponse(BaseModel):
    """Model for MCP JSON-RPC responses"""
    jsonrpc: str = "2.0"
    id: str
    result: Optional[Dict[str, Any]] = None
    error: Optional[Dict[str, Any]] = None


class ToolCallResult(BaseModel):
    """Model for tool call results"""
    name: str
    args: Dict[str, Any]
    result: Optional[str] = None
    error: Optional[str] = None


class MCPTool(BaseTool):
    """Wrapper for MCP server tools to integrate with LangChain"""
    
    name: str
    description: str
    model_config = ConfigDict(arbitrary_types_allowed=True)
    
    def __init__(self, name: str, description: str, mcp_server_process: subprocess.Popen, **kwargs: Any) -> None:
        super().__init__(name=name, description=description, **kwargs)
        self._mcp_server_process: subprocess.Popen[bytes] = mcp_server_process
    
    def _run(self, **kwargs: Any) -> str:
        """Execute the MCP tool with given parameters"""
        try:
            # Prepare MCP request
            request: MCPRequest = MCPRequest(
                id=f"tool_call_{self.name}",
                method="tools/call",
                params={
                    "name": self.name,
                    "arguments": kwargs
                }
            )
            
            # Send request to MCP server
            request_line: str = request.model_dump_json() + "\n"
            if self._mcp_server_process.stdin:
                self._mcp_server_process.stdin.write(request_line.encode())
                self._mcp_server_process.stdin.flush()
            
            # Read response
            if self._mcp_server_process.stdout:
                response_line: str = self._mcp_server_process.stdout.readline().decode().strip()
                response_data: Dict[str, Any] = json.loads(response_line)
                response: MCPResponse = MCPResponse(**response_data)
            else:
                return "Error: MCP server process stdout not available"
            
            if response.error:
                return f"Error: {response.error.get('message', 'Unknown error')}"
            
            # Extract content from response
            if response.result and "content" in response.result:
                content: List[Dict[str, Any]] = response.result["content"]
                if content and len(content) > 0:
                    return content[0].get("text", "No text content")
            
            return "No response from tool"
            
        except Exception as e:
            return f"Error executing tool {self.name}: {str(e)}"


class MCPClient:
    """Client for interacting with MCP servers"""
    
    def __init__(self, server_path: str) -> None:
        self.server_path: str = server_path
        self.server_process: Optional[subprocess.Popen[bytes]] = None
        self.tools: List[MCPTool] = []
    
    async def start_server(self) -> None:
        """Start the MCP server process"""
        try:
            self.server_process = subprocess.Popen(
                [self.server_path],
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=False  # We'll handle encoding manually
            )
            
            # Initialize the server
            init_params: Dict[str, Any] = {
                "protocolVersion": "2024-11-05",
                "capabilities": {},
                "clientInfo": {
                    "name": "llm-mcp-client-python",
                    "version": "1.0.0"
                }
            }
            await self._send_request("initialize", init_params)
            
            print(f"✓ Connected to MCP server: {self.server_path}")
            
        except Exception as e:
            raise Exception(f"Failed to start MCP server: {e}")
    
    async def _send_request(self, method: str, params: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Send a request to the MCP server"""
        if not self.server_process:
            raise Exception("MCP server not started")
        
        request: MCPRequest = MCPRequest(
            id=f"req_{method}",
            method=method,
            params=params or {}
        )
        
        try:
            request_line: str = request.model_dump_json() + "\n"
            if self.server_process and self.server_process.stdin:
                self.server_process.stdin.write(request_line.encode())
                self.server_process.stdin.flush()
            
            if self.server_process and self.server_process.stdout:
                response_line: str = self.server_process.stdout.readline().decode().strip()
                response_data: Dict[str, Any] = json.loads(response_line)
                response: MCPResponse = MCPResponse(**response_data)
            else:
                raise Exception("MCP server process not available")
            
            if response.error:
                raise Exception(f"MCP error: {response.error.get('message', 'Unknown error')}")
            
            return response.result or {}
            
        except Exception as e:
            raise Exception(f"Failed to communicate with MCP server: {e}")
    
    async def discover_tools(self) -> None:
        """Discover available tools from the MCP server"""
        try:
            result: Dict[str, Any] = await self._send_request("tools/list")
            tools_data: List[Dict[str, Any]] = result.get("tools", [])
            
            self.tools = []
            for tool_data in tools_data:
                if self.server_process:
                    tool: MCPTool = MCPTool(
                        name=tool_data["name"],
                        description=tool_data.get("description", ""),
                        mcp_server_process=self.server_process
                    )
                    self.tools.append(tool)
            
            print(f"✓ Discovered {len(self.tools)} tools: {[t.name for t in self.tools]}")
            
        except Exception as e:
            raise Exception(f"Failed to discover tools: {e}")
    
    def stop_server(self) -> None:
        """Stop the MCP server process"""
        if self.server_process:
            self.server_process.terminate()
            self.server_process.wait()


class LLMChat:
    """Chat interface with LLM and MCP tool integration"""
    
    def __init__(self, model_name: str, mcp_client: MCPClient) -> None:
        self.llm: ChatOllama = ChatOllama(model=model_name, temperature=0.7)
        self.mcp_client: MCPClient = mcp_client
        self.conversation_history: List[BaseMessage] = []
    
    def _create_system_message(self) -> SystemMessage:
        """Create system message with tool descriptions"""
        if not self.mcp_client.tools:
            return SystemMessage(content="You are a helpful AI assistant.")
        
        tool_descriptions = []
        for tool in self.mcp_client.tools:
            tool_descriptions.append(f"- {tool.name}: {tool.description}")
        
        tools_text = "\n".join(tool_descriptions)
        
        system_content = f"""You are a helpful AI assistant with access to the following tools:

{tools_text}

When you want to use a tool, respond with a special format:
USE_TOOL: tool_name(parameter1=value1, parameter2=value2)

For example:
- USE_TOOL: get_random_number(min=1, max=10)
- USE_TOOL: get_random_string(length=8, charset=alpha)
- USE_TOOL: get_random_choice(choices=["red", "blue", "green"])

You can use multiple tools in one response if needed. Be helpful and use tools when they would be useful to answer the user's question."""
        
        return SystemMessage(content=system_content)
    
    def _parse_tool_calls(self, response: str) -> List[Dict[str, Any]]:
        """Parse tool calls from LLM response"""
        
        # Find all tool calls in the response
        pattern: str = r'USE_TOOL:\s*(\w+)\((.*?)\)'
        matches: List[Tuple[str, str]] = re.findall(pattern, response)
        
        tool_calls: List[Dict[str, Any]] = []
        for tool_name, args_str in matches:
            # Parse arguments
            args: Dict[str, Any] = {}
            if args_str.strip():
                # Simple argument parsing for key=value pairs
                arg_pairs: List[str] = args_str.split(',')
                for pair in arg_pairs:
                    if '=' in pair:
                        key: str
                        value: Union[str, int, float, List[str]]
                        key, value_str = pair.split('=', 1)
                        key = key.strip()
                        value_str = value_str.strip()
                        
                        # Remove quotes if present
                        if value_str.startswith('"') and value_str.endswith('"'):
                            value = value_str[1:-1]
                        elif value_str.startswith('[') and value_str.endswith(']'):
                            # Handle array values
                            array_content: str = value_str[1:-1]
                            choices: List[str] = [choice.strip().strip('"') for choice in array_content.split(',')]
                            value = choices
                        else:
                            # Try to convert to number
                            try:
                                if '.' in value_str:
                                    value = float(value_str)
                                else:
                                    value = int(value_str)
                            except ValueError:
                                value = value_str  # Keep as string
                        
                        args[key] = value
            
            tool_calls.append({"name": tool_name, "args": args})
        
        return tool_calls
    
    async def chat(self, user_input: str) -> str:
        """Process user input and return response"""
        # Add user message to conversation
        self.conversation_history.append(HumanMessage(content=user_input))
        
        # Create messages for LLM
        messages = [self._create_system_message()] + self.conversation_history
        
        # Get initial response from LLM
        response = await self.llm.ainvoke(messages)
        response_text = response.content
        
        # Check if LLM wants to use tools
        tool_calls = self._parse_tool_calls(response_text)
        
        if tool_calls:
            print("🔧 Using tools...")
            
            tool_results = []
            for tool_call in tool_calls:
                tool_name = tool_call["name"]
                tool_args = tool_call["args"]
                
                # Find the tool
                tool = next((t for t in self.mcp_client.tools if t.name == tool_name), None)
                if tool:
                    try:
                        result = tool._run(**tool_args)
                        tool_results.append(f"{tool_name}: {result}")
                        print(f"  ✓ {tool_name} → {result}")
                    except Exception as e:
                        error_msg = f"Error using {tool_name}: {str(e)}"
                        tool_results.append(error_msg)
                        print(f"  ✗ {error_msg}")
                else:
                    error_msg = f"Tool {tool_name} not found"
                    tool_results.append(error_msg)
                    print(f"  ✗ {error_msg}")
            
            # Get final response with tool results
            tool_results_text = "\n".join(tool_results)
            follow_up_message = f"Tool results:\n{tool_results_text}\n\nBased on these results, provide a helpful response to the user."
            
            messages.append(AIMessage(content=response_text))
            messages.append(HumanMessage(content=follow_up_message))
            
            final_response = await self.llm.ainvoke(messages)
            final_text = final_response.content
            
            # Add to conversation history
            self.conversation_history.append(AIMessage(content=final_text))
            
            return final_text
        else:
            # No tools used, return direct response
            self.conversation_history.append(AIMessage(content=response_text))
            return response_text


async def main() -> int:
    """Main function"""
    parser: argparse.ArgumentParser = argparse.ArgumentParser(description="LLM MCP Client - Python implementation")
    parser.add_argument("mcp_server_path", help="Path to the MCP server executable")
    parser.add_argument("--model", default="gemma3:12b", help="Ollama model name (default: gemma3:12b)")
    
    args: argparse.Namespace = parser.parse_args()
    
    print("🚀 Starting LLM MCP Client (Python + LangChain)")
    print(f"📡 MCP Server: {args.mcp_server_path}")
    print(f"🤖 Model: {args.model}")
    print()
    
    # Initialize MCP client
    mcp_client: MCPClient = MCPClient(args.mcp_server_path)
    
    try:
        # Start MCP server and discover tools
        await mcp_client.start_server()
        await mcp_client.discover_tools()
        
        # Initialize chat
        chat: LLMChat = LLMChat(args.model, mcp_client)
        
        print()
        print("💬 Chat started! Type 'quit' or 'exit' to stop.")
        print()
        
        # Interactive chat loop
        while True:
            try:
                user_input: str = input("User: ").strip()
                
                if not user_input:
                    continue
                
                if user_input.lower() in ["quit", "exit"]:
                    break
                
                print()
                response: str = await chat.chat(user_input)
                print(f"Assistant: {response}")
                print()
                
            except KeyboardInterrupt:
                print("\n👋 Goodbye!")
                break
            except Exception as e:
                print(f"❌ Error: {e}")
                print()
    
    except Exception as e:
        print(f"❌ Failed to start: {e}")
        return 1
    
    finally:
        mcp_client.stop_server()
    
    return 0


if __name__ == "__main__":
    sys.exit(asyncio.run(main()))