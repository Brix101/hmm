use hmm_utils::reader::Reader;

fn main() {
        let file_list = Reader::list_files();

    for file in file_list {
        println!("{:#?}", file);
    }
}
