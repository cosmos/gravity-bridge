//! Main entry point for Gorc

#![deny(warnings, missing_docs, trivial_casts, unused_qualifications)]
#![forbid(unsafe_code)]

use gorc::application::APP;

/// Boot Gorc
fn main() {
    openssl_probe::init_ssl_cert_env_vars();
    abscissa_core::boot(&APP);
}
