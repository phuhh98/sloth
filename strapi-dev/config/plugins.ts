import type { Core } from "@strapi/strapi";
import NodeMailerProvider from "@strapi/provider-email-nodemailer";

const config = ({
  env,
}: Core.Config.Shared.ConfigParams): Core.Config.Plugin => ({
  // Upload provider https://docs.strapi.io/cloud/advanced/upload
  upload: {
    config: {
      providerOptions: {
        localServer: {
          maxage: 300000,
        },
      },
    },
  },

  // Email config https://docs.strapi.io/cms/features/email
  // Node mailer provider https://market.strapi.io/providers/@strapi-provider-email-nodemailer
  email: {
    config: {
      provider: "nodemailer",
      providerOptions: {
        host: env("SMTP_HOST", "smtp.example.com"),
        port: env("SMTP_PORT", 587),
        secure: false, // Use `true` for port 465
        ...(env("NODE_ENV") === "production" &&
        env("SMTP_USERNAME") &&
        env("SMTP_PASSWORD")
          ? {
              auth: {
                user: env("SMTP_USERNAME"),
                pass: env("SMTP_PASSWORD"),
              },
            }
          : {}),
        // ... any custom nodemailer options
      } satisfies Parameters<typeof NodeMailerProvider.init>[0],
      settings: {
        /**
         * TODO: Update these default email settings as needed.
         */
        defaultFrom: "hello@example.com",
        defaultReplyTo: "hello@example.com",
      } satisfies Parameters<typeof NodeMailerProvider.init>[1],
    },
  },

  "cheap-strapi-plugin": {
    enabled: true,
    // resolve: "../../../../cheap-strapi-plugin",
  },
});

export default config;
