# Agent Skills Documentation

本文档定义了 Krag 系统暴露的核心能力（Skills），旨在供其他 AI Agent（如基于 ReAct 的智能体、GPTs 等）通过工具调用（Function Calling）的方式集成和使用。

## 1. 知识库检索 (Knowledge Retrieval)

**Skill Name**: `search_knowledge_base`

**Description**:
当用户询问特定领域知识、私有文档内容或需要引用外部资料时，使用此工具。该工具会在向量数据库中进行语义检索，返回最相关的文档片段。

**Endpoint**: `POST /v1/knowledge/search`

**Parameters (JSON Schema)**:
```json
{
  "type": "object",
  "properties": {
    "query": {
      "type": "string",
      "description": "The search query based on user's question. Should be specific and contain keywords."
    },
    "k": {
      "type": "integer",
      "description": "Number of results to return. Default is 3.",
      "default": 3
    }
  },
  "required": ["query"]
}
```

**Example Call**:
```json
{
  "name": "search_knowledge_base",
  "arguments": {
    "query": "Krag 项目的架构设计是怎样的？",
    "k": 5
  }
}
```

---

## 2. 知识库上传 (Document Upload)

**Skill Name**: `upload_document`

**Description**:
当用户提供文件（PDF, TXT, Markdown）并希望将其加入知识库以便后续检索时，使用此工具。注意：此工具通常需要文件上传的前置步骤，适合支持文件处理的 Agent。

**Endpoint**: `POST /v1/knowledge/upload` (Multipart Form)

**Parameters**:
*   `file`: The file binary data.

---

## 3. 智能对话 (Chat)

**Skill Name**: `chat_completion`

**Description**:
调用 Krag 的核心 LLM 进行对话。支持上下文记忆。

**Endpoint**: `POST /v1/chat`

**Parameters (JSON Schema)**:
```json
{
  "type": "object",
  "properties": {
    "content": {
      "type": "string",
      "description": "The user's message or prompt."
    },
    "conversation_id": {
      "type": "string",
      "description": "Unique identifier for the conversation context. If empty, a new conversation is started."
    },
    "use_rag": {
      "type": "boolean",
      "description": "Whether to enable RAG (Retrieval Augmented Generation) for this message.",
      "default": false
    },
    "stream": {
      "type": "boolean",
      "description": "Whether to stream the response.",
      "default": false
    }
  },
  "required": ["content"]
}
```

---

## 集成指南

### Python (LangChain)
```python
from langchain.tools import Tool

def search_krag(query: str):
    # Call Krag API
    pass

tools = [
    Tool(
        name="search_knowledge_base",
        func=search_krag,
        description="Search for relevant documents in Krag knowledge base."
    )
]
```

### OpenAI Function Calling
直接将上述 `Parameters` 中的 JSON Schema 作为 `functions` 参数传递给 OpenAI API。
