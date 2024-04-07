import { z } from "zod";

export const notificationConfigFormValidator = z.object({
  mention: z.string().regex(/^[a-z0-9-]*$/, {
    message: "Lowercase letters, numbers, and “-” only.",
  }),
  statuses: z
    .object({
      status: z.string(),
    })
    .array(),
});
export type NotificationConfigFormData = z.infer<
  typeof notificationConfigFormValidator
>;

export const emptyNotificationConfig: NotificationConfigFormData = {
  mention: "",
  statuses: [],
};
