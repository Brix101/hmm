import { FileEntity } from "@/types/file.type";
import { STATIC_URL } from "@/constant/server.constant";
import { useBoundStore } from "@/store";
import { cn } from "@/lib/utils";
import { Folder } from "lucide-react";

interface Props {
  file: FileEntity;
}

const FileCard = ({ file }: Props) => {
  const { setActiveFilePath } = useBoundStore();
  const isFolder = file.files;
  const isImage = file.fileType?.includes("image");

  const handleNavigate = () => {
    if (isFolder) {
      setActiveFilePath(file.path);
    }
  };
  return (
    <div
      className={cn(
        "h-52 flex flex-col items-center border border-gray-50 justify-center truncate",
        isFolder ? "cursor-pointer hover:bg-blue-50" : ""
      )}
      onClick={handleNavigate}
    >
      {/* {file.fileType?.includes("image") ? ( */}
      {/*   <img src={STATIC_URL + file.path} alt={file.name} /> */}
      {/* ) : ( */}
      {/*   <CardFooter>{file.name}</CardFooter> */}
      {/* )} */}
      {isImage ? (
        <img src={STATIC_URL + file.path} alt={file.name} />
      ) : (
        <Folder size={170} />
      )}
      <span className="truncate">{file.name}</span>
    </div>
  );
};
export default FileCard;
