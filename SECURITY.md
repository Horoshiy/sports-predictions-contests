# Security Configuration Guide

## Environment Variables

### Critical Security Settings

**NEVER commit the `.env` file to version control!**

### Required Configuration for Production

1. **JWT Secret**
   ```bash
   # Generate a secure random string (at least 32 characters)
   JWT_SECRET=$(openssl rand -base64 32)
   ```

2. **Database Password**
   ```bash
   # Use a strong password with at least 16 characters
   DB_PASSWORD=<secure-random-password>
   ```

3. **SSL/TLS for Database**
   ```bash
   # ALWAYS use SSL in production
   DB_SSLMODE=require  # or verify-full for maximum security
   ```

### Development vs Production

| Setting | Development | Production |
|---------|-------------|------------|
| `JWT_SECRET` | Can use default | **MUST** be unique and secure |
| `DB_PASSWORD` | Can use default | **MUST** be strong and unique |
| `DB_SSLMODE` | `disable` (OK) | `require` or `verify-full` |
| `CORS_ALLOWED_ORIGINS` | `*` (OK) | Specific domains only |

## Password Requirements

The platform enforces the following password requirements:

- Minimum 8 characters
- At least one uppercase letter (A-Z)
- At least one lowercase letter (a-z)
- At least one number (0-9)

**Recommended**: 12+ characters with special characters for strong passwords.

## Email Validation

Email addresses must match the pattern:
- Valid characters before @: letters, numbers, dots, underscores, percent, plus, hyphen
- Valid domain with at least 2-character TLD

## Setup Instructions

### First Time Setup

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Generate secure secrets:
   ```bash
   # Generate JWT secret
   echo "JWT_SECRET=$(openssl rand -base64 32)" >> .env
   
   # Generate database password
   echo "DB_PASSWORD=$(openssl rand -base64 24)" >> .env
   ```

3. Configure for your environment:
   ```bash
   # For production, enable SSL
   echo "DB_SSLMODE=require" >> .env
   
   # Set allowed CORS origins
   echo "CORS_ALLOWED_ORIGINS=https://yourdomain.com" >> .env
   ```

4. Verify `.env` is in `.gitignore`:
   ```bash
   grep -q "^\.env$" .gitignore || echo ".env" >> .gitignore
   ```

### Production Deployment Checklist

- [ ] Unique JWT_SECRET generated
- [ ] Strong DB_PASSWORD set
- [ ] DB_SSLMODE set to `require` or `verify-full`
- [ ] CORS_ALLOWED_ORIGINS restricted to your domains
- [ ] `.env` file NOT committed to version control
- [ ] Secrets stored in secure secrets manager (AWS Secrets Manager, etc.)
- [ ] Regular security audits scheduled
- [ ] HTTPS enabled for all services
- [ ] Rate limiting configured on API Gateway

## Secrets Management

For production deployments, consider using:

- **AWS Secrets Manager**: Store and rotate secrets automatically
- **HashiCorp Vault**: Centralized secrets management
- **Kubernetes Secrets**: For K8s deployments
- **Docker Secrets**: For Docker Swarm deployments

## Security Best Practices

1. **Rotate secrets regularly** (every 90 days minimum)
2. **Use different secrets** for each environment (dev, staging, prod)
3. **Never log secrets** or include them in error messages
4. **Limit access** to production secrets to essential personnel only
5. **Monitor for exposed secrets** using tools like git-secrets or truffleHog
6. **Enable audit logging** for all secret access

## Reporting Security Issues

If you discover a security vulnerability, please email security@yourdomain.com instead of creating a public issue.

## Additional Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [PostgreSQL SSL Documentation](https://www.postgresql.org/docs/current/ssl-tcp.html)
