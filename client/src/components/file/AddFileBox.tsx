import {
  Sheet,
  SheetTrigger,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetFooter,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import Icons from "@/components/Icons";
import { Separator } from "@/components/ui/separator";
import { ScrollArea } from "@/components/ui/scroll-area";
import { FileWithPath, useDropzone } from "react-dropzone";
import { useCallback } from "react";
import { cn } from "@/lib/utils";

function AddFileBox() {
  const onDrop = useCallback((acceptedFiles: FileWithPath[]) => {
    // const file = acceptedFiles;
    console.log({ acceptedFiles });
  }, []);
  const {
    getRootProps,
    getInputProps,
    isDragActive,
    // isDragAccept,
    // isDragReject,
  } = useDropzone({
    onDrop,
    multiple: true,
    disabled: false,
  });

  return (
    <Sheet>
      <SheetTrigger asChild>
        <Card className="fixed right-10 bottom-10">
          <CardContent className="p-1">
            <Button variant={"outline"} size={"icon"}>
              <Icons.add />
            </Button>
          </CardContent>
        </Card>
      </SheetTrigger>
      <SheetContent className="flex flex-col pr-0 w-full sm:max-w-lg">
        <SheetHeader className="px-1">
          <SheetTitle>Upload File</SheetTitle>
        </SheetHeader>
        <Separator />

        <div className="flex justify-center items-center w-full">
          <div
            {...getRootProps({
              className: cn(
                "flex flex-col justify-center items-center mr-5 w-full h-32 rounded-lg border-2 border-dashed cursor-pointer  hover:bg-gray-100 ",
                isDragActive
                  ? "border-blue-500 bg-blue-50"
                  : "border-gray-300 bg-gray-50"
              ),
            })}
          >
            <div className="flex flex-col justify-center items-center pt-5 pb-6">
              <svg
                className="mb-4 w-8 h-8 text-gray-500 dark:text-gray-400"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 20 16"
              >
                <path
                  stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13 13h3a3 3 0 0 0 0-6h-.025A5.56 5.56 0 0 0 16 6.5 5.5 5.5 0 0 0 5.207 5.021C5.137 5.017 5.071 5 5 5a4 4 0 0 0 0 8h2.167M10 15V6m0 0L8 8m2-2 2 2"
                />
              </svg>
              <p className="mb-2 text-sm text-gray-500 dark:text-gray-400">
                <span className="font-semibold">Click to upload</span> or drag
                and drop
              </p>
            </div>
            <input className="hidden" {...getInputProps()} />
          </div>
        </div>
        <Separator />
        <ScrollArea className="h-full">
          <div className="flex flex-col gap-5 pr-6">items here</div>
        </ScrollArea>
        <div className="grid gap-1.5 pr-6 text-sm">
          <Separator className="mt-2" />

          <SheetFooter className="mt-1.5">
            <Button
              aria-label="Proceed to checkout"
              size="sm"
              className="w-full"
            >
              Upload Items
            </Button>
          </SheetFooter>
        </div>
      </SheetContent>
    </Sheet>
  );
}
export default AddFileBox;
