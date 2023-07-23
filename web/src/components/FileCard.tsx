import { FileEntity } from "@/types/file.type";
import { STATIC_URL } from "@/constant/server.constant";
import { useBoundStore } from "@/store";
import { cn } from "@/lib/utils";
import { Folder, File } from "lucide-react";

interface Props {
  file: FileEntity;
}

const FileCard = ({ file }: Props) => {
  const { setPathHistory } = useBoundStore();
  const isDir = file.isDir;
  const isImage = file.fileType?.includes("image");

  const handleNavigate = () => {
    if (isDir) {
      setPathHistory(file.path);
    } else {
      window.open(STATIC_URL + file.path, "_blank");
    }
  };
  return (
    <div
      className={cn(
        "h-52 flex flex-col items-center justify-center truncate",
        /*  isDir ?  */ "cursor-pointer hover:bg-blue-50" /* : "" */
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
      ) : isDir ? (
        <Folder size={170} />
      ) : (
        <File size={170} />
      )}
      <span className="truncate">{file.name}</span>
    </div>
  );
};
export default FileCard;
