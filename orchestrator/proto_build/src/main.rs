//! Protobuf files in the gravity repo, copying the result to the gravity_proto crate for import
//! and use. While this builder generates about a dozen files only one contains all the gravity
//! proto info and the rest are discarded in favor of upstream cosmos-sdk-proto

// Building new Gravity rust proto definitions
// run 'cargo run'
// go to gravity_proto/prost
// delete all files except gravity.v1.rs
// re-write calls to super::super::cosmos as cosmos-sdk-proto::cosmos

use std::path::Path;
use std::path::PathBuf;
use walkdir::WalkDir;

// All paths must end with a / and either be absolute or include a ./ to reference the current
// working directory.
/// A temporary directory for proto building

fn main() {
    let out_path: PathBuf = "../gravity_proto/src/prost/".parse().unwrap();
    compile_protos(&out_path);
}

fn compile_protos(out_dir: &Path) {
    println!(
        "[info] Compiling .proto files to Rust into '{}'...",
        out_dir.display()
    );

    let root = env!("CARGO_MANIFEST_DIR");
    let root: PathBuf = root.parse().unwrap();
    // this gives us the repo root by going up two levels from the module root
    let root = root.parent().unwrap().parent().unwrap().to_path_buf();

    let mut gravity_proto_dir = root.clone();
    gravity_proto_dir.push("module/proto/gravity/v1");
    let mut gravity_proto_include_dir = root.clone();
    gravity_proto_include_dir.push("module/proto");
    let mut third_party_proto_include_dir = root;
    third_party_proto_include_dir.push("module/third_party/proto");

    // Paths
    let proto_paths = [gravity_proto_dir];
    // we need to have an include which is just the folder of our protos to satisfy protoc
    // which insists that any passed file be included in a directory passed as an include
    let proto_include_paths = [gravity_proto_include_dir, third_party_proto_include_dir];

    // List available proto files
    let mut protos: Vec<PathBuf> = vec![];
    for proto_path in &proto_paths {
        protos.append(
            &mut WalkDir::new(proto_path)
                .into_iter()
                .filter_map(|e| e.ok())
                .filter(|e| {
                    e.file_type().is_file()
                        && e.path().extension().is_some()
                        && e.path().extension().unwrap() == "proto"
                })
                .map(|e| e.into_path())
                .collect(),
        );
    }

    // Compile all proto files
    let mut config = prost_build::Config::default();
    config.out_dir(out_dir);
    config
        .compile_protos(&protos, &proto_include_paths)
        .unwrap();

    // Compile all proto client for GRPC services
    println!("[info ] Compiling proto clients for GRPC services!");
    tonic_build::configure()
        .build_client(true)
        .build_server(false)
        .format(false)
        .out_dir(out_dir)
        .compile(&protos, &proto_include_paths)
        .unwrap();

    println!("[info ] => Done!");
}
