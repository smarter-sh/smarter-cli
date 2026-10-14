# chat

Chat with an LLMClient from the command line, through the platform's cli
prompt api (`/api/v1/cli/prompt/<name>/` and `/api/v1/cli/prompt/config/<name>/`).

```console
smarter chat prompt --llmclient <name> --prompt "Hello"
smarter chat config --llmclient <name>
```

`--chatbot` is a deprecated alias of `--llmclient`: LLMClient replaced ChatBot
in v0.14.
