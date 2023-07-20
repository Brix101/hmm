import { Sheet, SheetTrigger, SheetContent } from "../ui/sheet";
import { Button } from "../ui/button";

function AddFileBox() {
  return (
    <Sheet>
      <SheetTrigger>
        <Button
          className="fixed right-10 bottom-10 shadow-lg"
          variant={"default"}
        >
          Add File
        </Button>
      </SheetTrigger>
      <SheetContent className="flex flex-col pr-0 w-full sm:max-w-lg"></SheetContent>
    </Sheet>
  );
}
export default AddFileBox;
