import z from "zod";

const baseFileSchema = z.object({
  name: z.string(),
  size: z.number(),
  path: z.string(),
  fileType: z.string().optional().nullish(),
});

export type FileEntity = z.infer<typeof baseFileSchema> & {
  files?: FileEntity[];
};

export const fileEntitySchema: z.ZodType<FileEntity> = baseFileSchema.extend({
  files: z.lazy(() => fileEntitySchema.array().optional()),
});
