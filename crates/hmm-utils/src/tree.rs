use glob::glob;
use serde::{Deserialize, Serialize};
use std::{cmp::Ordering, fs, path::PathBuf, time::SystemTime};

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
pub struct Explorer {
    path: String,
    hidden: bool,
}

impl Explorer {
    pub fn new() -> Self {
        let home_dir = dirs::home_dir().expect("Failed to get home directory");
        let path = home_dir.to_string_lossy().to_string();
        let hidden = false;

        Explorer { path, hidden }
    }

    pub fn explore(&self, path: Option<String>, hidden: Option<bool>) -> Vec<FileInfo> {
        let view_hidden = hidden.unwrap_or(self.hidden);
        let root_path = path.unwrap_or(self.path.clone());

        let pattern = format!("{}/*", root_path);

        let mut files = glob(&pattern)
            .expect("Failed to read glob pattern")
            .filter_map(|entry| {
                if let Ok(path) = entry {
                    let metadata = fs::metadata(&path).expect("Failed to read metadata");
                    let file_type = metadata.file_type();
                    let name = path.file_name().unwrap().to_string_lossy().to_string();

                    let is_hidden = name.starts_with(".");

                    if view_hidden || !is_hidden {
                        Some(FileInfo {
                            name,
                            path: path.clone(),
                            size: metadata.len(),
                            file_ext: path
                                .extension()
                                .and_then(|os_str| os_str.to_str())
                                .map(|s| s.to_string()),
                            is_dir: file_type.is_dir(),
                            mod_time: metadata.modified().unwrap_or(SystemTime::UNIX_EPOCH),
                        })
                    } else {
                        None
                    }
                } else {
                    None
                }
            })
            .collect::<Vec<_>>();

        files.sort_by(|a, b| {
            let is_dir_cmp = b.is_dir.cmp(&a.is_dir); // Compare in reverse order for directories to appear first
            if is_dir_cmp == Ordering::Equal {
                a.name.cmp(&b.name)
            } else {
                is_dir_cmp
            }
        });

        return files;
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_explore() {
        let explorer = Explorer::new();
        let files = explorer.explore(None, None);

        assert!(!files.is_empty(), "No files were viewed");
    }

    #[test]
    fn test_explore_invalid_path() {
        let explorer = Explorer::new();
        let files_invalid_path =
            explorer.explore(Some("/non_existent_directory".to_string()), None);

        assert!(
            files_invalid_path.is_empty(),
            "Invalid path should result in an empty view"
        );
    }

    #[test]
    fn test_explore_hidden_false() {
        let explorer = Explorer::new();
        let files_hidden_false = explorer.explore(None, Some(false));
        assert!(
            !files_hidden_false
                .iter()
                .any(|file| file.name.starts_with('.')),
            "Hidden files should not be viewed"
        );
    }

    #[test]
    fn test_explore_hidden_true() {
        let explorer = Explorer::new();
        let files_hidden_true = explorer.explore(None, Some(true));
        assert!(
            !files_hidden_true.is_empty(),
            "No files were viewed when 'hidden' is true"
        );
    }
}
