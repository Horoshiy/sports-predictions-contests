# Sports Prediction Contests Documentation

[🇺🇸 English](en/) | [🇷🇺 Русский](ru/)

## Overview

Welcome to the comprehensive documentation for the Sports Prediction Contests platform - a multilingual, multi-sport API-first platform for creating and running sports prediction competitions.

## Quick Links

### English Documentation
- [📚 Complete Documentation](en/README.md)
- [🚀 Quick Start Guide](en/deployment/quick-start.md)
- [📖 API Reference](en/api/services-overview.md)
- [🧪 Testing Guide](en/testing/e2e-testing.md)
- [🔧 Troubleshooting](en/troubleshooting/common-issues.md)

### Русская документация
- [📚 Полная документация](ru/README.md)
- [🚀 Быстрый старт](ru/deployment/quick-start.md)
- [📖 Справочник API](ru/api/services-overview.md)
- [🧪 Руководство по тестированию](ru/testing/e2e-testing.md)
- [🔧 Устранение неполадок](ru/troubleshooting/common-issues.md)

## Platform Features

- **🏆 Contest Constructor**: Customizable rules, scoring systems, and sport types
- **🌐 Multi-platform Support**: Web, mobile apps, Telegram/Facebook bots
- **⚡ API-First Architecture**: gRPC-based microservices with open API
- **📊 Real-time Updates**: Live scoring and leaderboards
- **🎮 Gamification**: Statistics tracking, achievements, and rankings

## Architecture

The platform consists of 7 core microservices:

| Service | Port | Purpose |
|---------|------|---------|
| API Gateway | 8080 | HTTP REST entry point |
| User Service | 8084 | Authentication & user management |
| Contest Service | 8085 | Contest & team management |
| Prediction Service | 8086 | Predictions & events |
| Scoring Service | 8087 | Scoring & leaderboards |
| Sports Service | 8088 | Sports data & sync |
| Notification Service | 8089 | Multi-channel notifications |

## Getting Started

Choose your preferred language and follow the quick start guide:

- **English**: [Quick Start Guide](en/deployment/quick-start.md)
- **Русский**: [Руководство по быстрому старту](ru/deployment/quick-start.md)

## Support

For issues and questions:
- Check the [troubleshooting guides](en/troubleshooting/common-issues.md)
- Review the [API documentation](en/api/services-overview.md)
- Examine the [testing procedures](en/testing/e2e-testing.md)

---

*Built for the Dynamous Kiro Hackathon - January 2026*
