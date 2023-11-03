use std::{fs, path::PathBuf, time::SystemTime};
use glob::glob;
use serde::{Deserialize, Serialize};


#[derive(Debug, Deserialize, Serialize)]
pub struct FileInfo {
    pub path: PathBuf,
    pub name: String,
    pub size: u64,
    pub file_ext: Option<String>,
    pub is_dir: bool,
    pub mod_time: SystemTime,
}

#[derive(Debug, Deserialize, Serialize)]
pub struct ViewFilesOptions {
    pub path: Option<String>,
    pub hidden: bool,
}

pub struct Reader;

impl Reader {
    pub fn view_files(args:ViewFilesOptions) -> Vec<FileInfo> {
        let hidden = args.hidden;

        let home_dir = dirs::home_dir().expect("Failed to get home directory");
        let root_path = args.path.unwrap_or(home_dir.to_string_lossy().to_string());
        let pattern = format!("{}/*", root_path);

        glob(&pattern)
            .expect("Failed to read glob pattern")
            .filter_map(|entry|{
                if let Ok(path) = entry {
                    let metadata = fs::metadata(&path).expect("Failed to read metadata");
                    let file_type = metadata.file_type();
                    let name= path.file_name().unwrap().to_string_lossy().to_string();

                    let is_hidden = name.starts_with(".");

                    if hidden || !is_hidden {
                        Some(FileInfo {
                            name,
                            path:path.clone(),
                            size:metadata.len(),
                            file_ext:path.extension().and_then(|os_str| os_str.to_str()).map(|s| s.to_string()),
                            is_dir:file_type.is_dir(),
                            mod_time:metadata.modified().unwrap_or(SystemTime::UNIX_EPOCH),
                        })
                    }else{
                        None
                    }
                } else {
                    None
                }
            }).collect()

    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_list_files() {
        let files = Reader::view_files(ViewFilesOptions { path: None, hidden: true });
        assert!(!files.is_empty(), "No files were listed");
    }
}
