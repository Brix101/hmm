import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { STATIC_URL } from "@/constant/server.constant";
import { useQueryFile } from "@/services/file.service";
import { useBoundStore } from "@/store";
import { Fragment } from "react";

const Home = () => {
  const {
    file: { fileUrls },
    appendToUrl,
    resetUrl,
  } = useBoundStore();

  const { data, isLoading } = useQueryFile(fileUrls ?? "");
  if (isLoading) {
    return <h1>Loading</h1>;
  }

  function handleResetClick() {
    resetUrl();
  }
  return (
    <>
      <div>
        <div className="flex items-center space-x-1 text-sm capitalize text-muted-foreground">
          <div
            className="cursor-pointer hover:underline truncate"
            onClick={handleResetClick}
          >
            Files
          </div>
          {fileUrls
            ?.replace("/", "")
            .split("/")
            .map((url, index) => (
              <Fragment key={index}>
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="w-4 h-4"
                  aria-hidden="true"
                >
                  <polyline points="9 18 15 12 9 6"></polyline>
                </svg>
                <div
                  className={
                    fileUrls.length - 1 === index ? "text-foreground" : ""
                  }
                >
                  {url}
                  {index}
                </div>
              </Fragment>
            ))}
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {data?.files?.map((file, index) => (
          <Card className="w-[350px] h-[350px]" key={file.name + index}>
            <CardContent>
              {file.fileType?.includes("image") ? (
                <img src={STATIC_URL + file.path} alt={file.name} />
              ) : (
                <>{file.name}</>
              )}
              {file.files ? (
                <Button onClick={() => appendToUrl(file.path)}>navigate</Button>
              ) : undefined}
            </CardContent>
          </Card>
        ))}
      </div>
    </>
  );
};

export default Home;
