#[macro_use]
extern crate error_chain;
#[macro_use]
extern crate log;
extern crate env_logger;
extern crate clap;
extern crate toml;
extern crate rustc_serialize;

mod config;
mod errors;

include!(concat!(env!("OUT_DIR"), "/version.rs"));

use std::path::Path;

use clap::{Arg, App};

fn main() {
    // Set up and parse the command-line arguments
    let matches = App::new("Ataxia Engine")
        .version(env!("CARGO_PKG_VERSION"))
        .author("Xenith Studios (see AUTHORS)")
        .about(env!("CARGO_PKG_DESCRIPTION"))
        .arg(Arg::with_name("config")
            .help("The config file to use")
            .short("c")
            .long("config")
            .value_name("FILE")
            .takes_value(true)
            .default_value("data/config.toml"))
        .arg(Arg::with_name("listen_addr")
            .help("Listen address and port")
            .short("l")
            .long("listen")
            .value_name("address:port")
            .takes_value(true))
        .arg(Arg::with_name("pid_file")
            .help("The filename to write the PID into")
            .short("p")
            .long("pid")
            .value_name("FILE")
            .takes_value(true))
        .arg(Arg::with_name("debug")
            .help("Enable debugging output")
            .short("d")
            .multiple(true))
        .arg(Arg::with_name("verbose")
            .help("Enable verbose output")
            .short("v")
            .multiple(true))
        .get_matches();

    // Initialize logging subsystem
    // TODO: Repalce with a propper logging system.
    env_logger::init().expect("Failed to initialize logging.");

    info!("Loading Ataxia Engine, compiled on {}", ATAXIA_COMPILED);

    // Load settings from config file
    let config_path = Path::new(matches.value_of("config")
        .expect("Unable to specify config file path."));
    info!("Loading configuration from: {:?}", config_path);
    let config = config::Config::read_config(config_path)
        .expect("Unable to load the configuration.");

    // Clean up from previous unclean shutdown if necessary
    //   Delete PID file if it exists

    // Set up callbacks for signals
    // Write PID file

    // Initialize
    //   Seed rand
    //   Environment
    //   Queues
    //   Database
    //   Lua

    // Initialize engine
    // Load initial game state
    //   Load database
    //   Load commands
    //   Load scripts
    //   Load world
    //   Load entities


    // Initialize networking event loop in dedicated thread
    // Spawn other threads?
    // Start main game loop

    // Shutdown is caught here?
    // Clean up
    //   Save the world
    //   Shutdown Lua
    //   Flush pending database writes
    //   Close database connection
    //   Remove PID file
}
