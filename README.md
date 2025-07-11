# 🍵 Teapot

> Bootstrap modern full-stack monorepos in seconds

Teapot is a powerful CLI tool that helps developers quickly scaffold production-ready monorepo projects with all the modern tooling you need. Whether you're building a web app, mobile app, or full-stack platform, Teapot sets up everything from your framework choices to your infrastructure, so you can focus on building your product.

## ✨ Features

- **🚀 Instant Setup** - Generate a complete Turborepo monorepo structure in seconds
- **📦 Modern Package Management** - Uses Bun for blazing-fast dependency management
- **🎯 Framework Flexibility** - Choose from React, Next.js, NestJS, Expo, or combine them all
- **🔧 Pre-configured Tooling** - ESLint, Prettier, and TypeScript configured out of the box
- **🐳 Infrastructure Ready** - Optional Docker Compose setup with Redis, Postgres, and more
- **☸️ Cloud Native** - Scaffold Kubernetes deployments with Pulumi
- **🤖 CI/CD Pipeline** - GitHub Actions workflows configured and ready to go
- **🎨 Beautiful CLI** - Clean, interactive interface powered by Bubble Tea

## 🛠️ Tech Stack

![Turborepo](https://img.shields.io/badge/Turborepo-000000?style=for-the-badge&logo=turborepo&logoColor=white)
![Bun](https://img.shields.io/badge/Bun-000000?style=for-the-badge&logo=bun&logoColor=white)
![React](https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB)
![Next.js](https://img.shields.io/badge/Next.js-000000?style=for-the-badge&logo=next.js&logoColor=white)
![NestJS](https://img.shields.io/badge/NestJS-E0234E?style=for-the-badge&logo=nestjs&logoColor=white)
![Expo](https://img.shields.io/badge/Expo-000020?style=for-the-badge&logo=expo&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Kubernetes](https://img.shields.io/badge/Kubernetes-326CE5?style=for-the-badge&logo=kubernetes&logoColor=white)
![Pulumi](https://img.shields.io/badge/Pulumi-8A3391?style=for-the-badge&logo=pulumi&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=github-actions&logoColor=white)

## 📦 Installation

Install Teapot globally using Go:

```bash
go install github.com/yourname/teapot@latest
```

Make sure your `$GOPATH/bin` is in your `PATH`.

## 🚀 Quick Start

1. **Create a new project:**
   ```bash
   teapot init my-awesome-project
   ```

2. **Follow the interactive prompts:**
   - Select your desired apps (React, Next.js, NestJS, Expo)
   - Choose infrastructure services (Redis, Postgres, etc.)
   - Configure deployment options (Docker Compose, Kubernetes)

3. **Navigate to your project and start developing:**
   ```bash
   cd my-awesome-project
   bun install
   bun dev
   ```

That's it! Your monorepo is ready with all the tooling configured.

## 🤝 Contributing

We love contributions! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please make sure to update tests as appropriate and follow the existing code style.

### Development Setup

```bash
# Clone the repo
git clone https://github.com/yourname/teapot.git
cd teapot

# Install dependencies
go mod tidy

# Run the CLI locally
go run main.go
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Made with ❤️ by the open source community