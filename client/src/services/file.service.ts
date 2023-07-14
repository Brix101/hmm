import { apiClient } from "@/lib/httpCommon";
import { FileEntity, fileEntitySchema } from "@/types/file.type";
import { UseQueryOptions, useQuery } from "@tanstack/react-query";
import { AxiosError } from "axios";

const getFiles = async () => {
  const res = await apiClient.get("/files");
  return fileEntitySchema.parse(res.data);
};

export const useGetFiles = (
  options?: UseQueryOptions<
    FileEntity,
    AxiosError,
    FileEntity,
    readonly [string]
  >
) => {
  return useQuery({
    queryKey: ["files"],
    queryFn: getFiles,
    ...options,
  });
};
