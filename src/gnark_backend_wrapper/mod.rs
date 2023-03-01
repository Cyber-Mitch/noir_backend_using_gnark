cfg_if::cfg_if! {
    if #[cfg(feature = "groth16")] {
        mod groth16;
        pub use groth16::{AddTerm, Fr, GoString, MulTerm, RawGate, RawR1CS};
        pub use groth16::verify_with_meta;
        pub use groth16::prove_with_meta;
        pub use groth16::verify_with_vk;
        pub use groth16::prove_with_pk;
        pub use groth16::get_exact_circuit_size;
        pub use groth16::preprocess;
    }
}
