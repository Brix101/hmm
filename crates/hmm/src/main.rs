use hmm_utils::explorer::Explorer;

fn main() {
    let explorer = Explorer::new();
    let files = explorer.explore(
        Some("/home/brix/Workspaces/brixterporras".to_string()),
        Some(true),
    );

    for file in files {
        println!("{:#?}", file);
    }
}
