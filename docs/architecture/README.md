# Architecture
Source adapter captures binding and generation at observation creation, signs an
opaque payload digest, and transports the entire envelope. Receiver supplies its
own trust policy and explicit destination expectation. Parser -> crypto/policy
filter -> accepted evidence set -> historical evaluator -> bounded result.

The core is pure batch computation. It has no persistence, telemetry, global
registry, network transport or clinical data parser. Integrators store/disclose
only necessary evidence and own transport encryption/access control. History union
reconciles partitions; consumers re-evaluate with current local policy. Old outputs
are snapshots, never durable permission or globally fresh identity assertions.
