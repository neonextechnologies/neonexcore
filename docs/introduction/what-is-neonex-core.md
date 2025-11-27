# What is Neonex Core?

Neonex Core is a modern, modular Go framework designed for building high-performance web applications and APIs. It combines the speed of Go with clean architecture patterns and developer-friendly tooling.

## Philosophy

Neonex Core is built on three core principles:

### 1. **Modularity First**

Applications are composed of self-contained modules that can be developed, tested, and deployed independently. Each module encapsulates its own:

* Models and business logic
* Routes and controllers
* Services and repositories
* Dependencies and configuration

### 2. **Developer Experience**

We believe that developer productivity is paramount. Neonex Core provides:

* **CLI Tools** - Scaffold projects and modules instantly
* **Hot Reload** - See changes immediately without manual restarts
* **Clear Structure** - Intuitive project organization
* **Rich Logging** - Understand what's happening in your application

### 3. **Performance Without Compromise**

Built on top of [Fiber](https://gofiber.io/) (which uses [fasthttp](https://github.com/valyala/fasthttp)), Neonex Core delivers exceptional performance while maintaining code clarity and maintainability.

## What Makes It Different?

### Compared to Other Go Frameworks

| Feature | Neonex Core | Gin | Echo | Fiber |
|---------|------------|-----|------|-------|
| Module System | ✅ Built-in | ❌ Manual | ❌ Manual | ❌ Manual |
| DI Container | ✅ Built-in | ❌ Manual | ❌ Manual | ❌ Manual |
| ORM Layer | ✅ Integrated | 🔄 Separate | 🔄 Separate | 🔄 Separate |
| CLI Tools | ✅ Full Suite | ❌ None | ❌ None | ❌ Basic |
| Hot Reload | ✅ Integrated | 🔄 Air | 🔄 Air | 🔄 Air |

### Architecture Approach

Neonex Core implements a **Modular Monolith** architecture:

```
┌─────────────────────────────────────────┐
│         Your Application                │
├─────────────────────────────────────────┤
│  Module 1  │  Module 2  │  Module 3   │
├────────────┼────────────┼──────────────┤
│         Neonex Core Layer              │
│  (DI, Routing, ORM, Logging, CLI)     │
├─────────────────────────────────────────┤
│         Fiber / fasthttp               │
└─────────────────────────────────────────┘
```

This approach provides:

* **Simplicity** - Single deployable unit
* **Modularity** - Clear boundaries between features
* **Scalability** - Can be split into microservices later if needed
* **Development Speed** - Faster than distributed systems

## Use Cases

Neonex Core is perfect for:

* ✅ **RESTful APIs** - Fast, well-structured API backends
* ✅ **Web Applications** - Server-side rendered or API-driven SPAs
* ✅ **Microservices** - Individual services in a distributed system
* ✅ **Internal Tools** - Admin panels, dashboards, automation tools
* ✅ **MVPs** - Rapid prototyping and validation

## Not Suitable For

* ❌ Simple static websites (use Hugo, Next.js instead)
* ❌ Real-time gaming servers (use dedicated game engines)
* ❌ Ultra-low latency trading systems (use C++/Rust)

## Next Steps

Ready to get started?

* [Install Neonex Core](../getting-started/installation.md)
* [Create your first project](../getting-started/quick-start.md)
* [Explore the architecture](architecture.md)
