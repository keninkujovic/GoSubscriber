# GoSubscriber Documentation

## Overview
GoSubscriber is a robust email collector built with Go, designed for straightforward integration into systems requiring email subscription functionalities. It's particularly useful for assembling email subscription lists, waiting lists, and update notification lists, all within a Dockerized environment for easy deployment and scalability.

## Key Features
- **Email Subscription Lists**: Efficiently collect emails for newsletters or promotional campaigns.
- **Waiting Lists**: Compile emails for upcoming product launches or services.
- **Update Lists**: Gather user emails to notify about new updates, releases, or content.

## Deployment
Deploy GoSubscriber quickly with Docker. Use the provided Docker Compose file for hassle-free setup, ensuring the application is ready to collect emails in no time.

## Security
GoSubscriber employs Basic Authentication for secure access to sensitive endpoints, ensuring that only authorized users can retrieve or modify email lists.

## Environment Variables
Configure GoSubscriber securely using environment variables for sensitive information, including database credentials and authentication details for protected endpoints.

- `BASIC_AUTH_USERNAME`: Username for Basic Authentication.
- `BASIC_AUTH_PASSWORD`: Password for Basic Authentication.

Ensure these variables are set in your environment to secure the application.

## Endpoints

### Subscribe
- **`/subscribe`**: Add an email to the collection by sending a POST request with the email in JSON format:
  ```json
  { "email": "example@email.com" }
  ```
  This endpoint is publicly accessible.

### Retrieve Emails (Protected)
- **`/emails`**: Fetch a list of all collected emails. This endpoint is protected with Basic Authentication and requires valid credentials:
  ```bash
  curl -u [username]:[password] http://yourdomain.com/emails
  ```

## Using GoSubscriber

### Adding an Email
To subscribe an email, send a POST request to `/subscribe` with the email address:

```bash
curl -X POST http://yourdomain.com/subscribe \
-H "Content-Type: application/json" \
-d '{"email": "user@example.com"}'
```

### Accessing Email List
To access the collected emails, use the `/emails` endpoint with Basic Authentication:

```bash
curl -u BASIC_AUTH_USERNAME:BASIC_AUTH_PASSWORD http://yourdomain.com/emails
```

Replace `BASIC_AUTH_USERNAME` and `BASIC_AUTH_PASSWORD` with your environment-configured credentials.

## Conclusion
GoSubscriber simplifies email collection and management, offering a Dockerized solution perfect for developers seeking an efficient way to integrate email subscription functionalities into their applications. With secure access controls and easy configuration through environment variables, GoSubscriber stands as a prime choice for managing email lists securely and efficiently.
