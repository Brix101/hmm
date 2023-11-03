use hmm_utils::reader::{Reader, ViewFilesOptions};

fn main() {
    let file_list = Reader::view_files(ViewFilesOptions{
        path:None,
        hidden:false
    });

    for file in file_list {
        println!("{:#?}", file);
    }
}
