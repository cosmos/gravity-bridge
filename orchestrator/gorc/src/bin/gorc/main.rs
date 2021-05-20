//! Main entry point for Gorc

#![deny(warnings, missing_docs, trivial_casts, unused_qualifications)]
#![forbid(unsafe_code)]

use gorc::application::APP;

/// Boot Gorc
fn main() {
    abscissa_core::boot(&APP);
}
