use hmm_utils::tree::Explorer;

fn main() {
    let explorer = Explorer::new();
    let files = explorer.explore(None, None);

    for file in files {
        println!("{:#?}", file);
    }
}
