use std::{fs, path::PathBuf, time::SystemTime};
use glob::glob;
use serde::{Deserialize, Serialize};


#[derive(Debug, Serialize, Deserialize)]
pub struct FileInfo {
    pub path: Option<PathBuf>,
    pub name: String,
    pub size: u64,
    pub file_ext: Option<String>,
    pub is_dir: bool,
    pub mod_time: SystemTime,
}


pub struct Reader;

impl Reader {
    pub fn list_files() -> Vec<FileInfo> {
        let home_dir = dirs::home_dir().expect("Failed to get home directory");
        let pattern = format!("{}/*", home_dir.to_string_lossy());

        glob(&pattern)
            .expect("Failed to read glob pattern")
            .filter_map(|entry| {
                if let Ok(path) = entry {
                    let metadata = fs::metadata(&path).ok()?;
                    let file_type = metadata.file_type();
                    let is_dir = file_type.is_dir();
                    let file_ext = path.extension()?.to_string_lossy().to_string();

                    Some(FileInfo {
                        path: Some(path.clone()),
                        name: path.file_name()?.to_string_lossy().to_string(),
                        size: metadata.len(),
                        file_ext: if !is_dir { Some(file_ext) } else { None },
                        is_dir,
                        mod_time: metadata.modified().unwrap_or(SystemTime::UNIX_EPOCH),
                    })
                } else {
                    None
                }
            })
            .collect()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_list_files() {
        let files = Reader::list_files();
        assert!(!files.is_empty(), "No files were listed");
    }
}
