# ADR-0001: MVP scope, execution boundaries, and data flow

- Status: Accepted for implementation
- Date: 2026-09-08
- Issue: #2
- Related: #1, #33, #34-#41

## Context

ARBEX aggregates liquidity from external LPs/exchanges/RFQ venues and DeFi into a synthetic book. The main architectural risk is treating heterogeneous execution venues as if they were atomically executable. The MVP therefore needs explicit boundaries for market data, inventory, routing, OMS, chain execution, recovery, accounting, and owner controls.

## Decision

### D01 Initial markets

The MVP supports a common venue model but does not assume a specific external LP is already approved. Initial implementation targets:

1. Base DEX connectivity for USDC/WETH using verified Uniswap V3 and Aerodrome pools.
2. One external spot venue adapter selected later from venues that satisfy the capability gate below.
3. Additional LP/RFQ adapters remain observe-only until capability, account eligibility, fee model, reconciliation endpoints, and recovery route are verified.

A venue is not live-enabled by configuration name alone.

### D02 Inventory placement

Use pre-funded inventory per `venue/account/asset`. No cross-venue transfer is part of an arbitrage transaction. Funds are reserved before plan creation and remain reserved through unknown/cancel-pending/chain-pending states. Bot credentials must not have withdrawal permission. Initial rebalancing is an Owner operation outside the arbitrage execution path.

### D03 DeFi pool types

Uniswap V3 and Aerodrome volatile pools are P0. Slipstream is P1 unless registry verification shows that the selected MVP pair lacks sufficient supported liquidity without it; in that case #14 is promoted into the execution dependency chain.

### D04 Success criteria

An opportunity is eligible only when expected net profit is positive after venue fees, DEX price impact, L2 gas, L1 data fee, expected recovery/rebalance cost, configured uncertainty buffer, and inventory constraints. Shadow PnL, expected PnL, realized trading PnL, recovery PnL, and fixed infrastructure cost must remain distinct metrics.

### D05 Connectivity and placement

Provider/venue selection is based on measured feed freshness, order/RPC round-trip latency, p50/p95/p99, disconnect behavior, rate limits, supported reconciliation APIs, region availability, and cost. Go process speed is not used as a substitute for network measurements.

### D06 UI priority

The first operational UI prioritizes synthetic book, source freshness, venue inventory, execution groups/child orders/fills, known/possible exposure, engine mode, and stop scope. It must never silently substitute demo data for missing production data.

### D07 Live limits

Live mode is fail-closed. Before activation the Owner must explicitly configure at least:

- allowed venues/accounts/assets/instruments;
- max notional per plan and per venue;
- max reserved capital;
- max known and possible exposure;
- max unhedged duration;
- normal trading loss/gas budget;
- separate recovery loss/gas budget;
- stale-data and quote-expiry limits;
- price deviation limits;
- recovery venues and maximum attempts;
- daily stop limits.

Missing limits keep the engine in observe/shadow mode.

### D08 Contract model

Use separate non-upgradeable contracts for atomic arbitrage and single-leg chain execution in the MVP. `AtomicArbitrageExecutor` enforces same-token end balance increase for atomic DEX round trips. `SingleLegExecutor` is for LP↔DEX recovery/second-leg execution and uses independent min-out, reference-price, role, budget, and pause controls. Arbitrary call/delegatecall is forbidden.

### D09 LP capability gate

A venue may progress beyond observe-only only if ARBEX can verify the capabilities needed for its configured execution mode. For non-atomic live trading the minimum gate is:

- deterministic instrument and asset mapping;
- authenticated market data with sequence/snapshot recovery or an equivalent consistency mechanism;
- client order ID or another idempotent reconciliation key;
- order lookup and fill/trade history;
- balance/free/hold information;
- cancel support with explicit state reconciliation;
- fee asset and fee calculation visibility;
- min quantity/notional and price/quantity increments;
- rate-limit/error semantics;
- API credentials restricted to read/trade/cancel and no withdrawal.

IOC/FOK is useful but not treated as atomicity. If an endpoint cannot distinguish unknown from not-found safely, the venue remains observe-only.

### D10 Execution order

Initial non-atomic execution is sequential, not simultaneous. The router chooses the first leg based on executable liquidity, expected slippage, inventory, cancellation behavior, and recoverability. The second leg quantity is based on actual net fill of the first leg. Simultaneous leg submission is disabled until failure replay demonstrates bounded exposure and a separate ADR approves it.

Atomic same-chain DEX↔DEX routes remain a separate execution class.

### D11 Unknown state and stop semantics

Unknown external side effects are reconciled before replacement or reservation release. Timeout, 404, stream loss, or cancel acknowledgement alone does not prove that an order did not fill.

Stop states are:

- `ENTRY_PAUSED`: no new arbitrage groups; reconciliation and approved recovery may continue.
- `RECOVERY_ONLY`: no new arbitrage groups; only reconciliation and recovery actions inside recovery policy/budget may execute.
- `HALT_ALL`: no new trading or recovery orders are created; existing external/chain side effects are still observed and reconciled.

Emergency hedging cannot bypass `HALT_ALL`.

### D12 Inventory rebalancing

Each venue/account/asset has target/min/max inventory bands. Rebalancing is not counted as arbitrage execution and its fees/PnL are accounted separately. The MVP only generates rebalance proposals and allows Owner-executed transfers. Automated withdrawals/transfers are out of scope.

## Execution boundary

The system is divided into the following responsibility chain:

`Market adapters -> normalized liquidity -> synthetic book -> inventory-aware router -> reservation -> OMS/execution -> reconciliation -> exposure/recovery -> ledger/API/UI`

External LP orders and chain transactions are separate side effects joined by an `ExecutionGroup`; they are never modeled as one atomic transaction.

## Data ownership

- Market adapters own raw venue/source observations and source revisions.
- Synthetic book owns normalized executable cost curves, not balances.
- Inventory owns settled/free/hold/pending/local reservations per venue/account/asset.
- Router owns immutable execution plans and their source/config revisions.
- OMS owns child-order lifecycle and fill reconciliation.
- Chain execution owns signer/nonce/transaction lifecycle.
- Recovery owns known/possible exposure and recovery actions.
- Ledger owns realized cash flows, fees, corrections, and reconciliation status.
- API/UI are projections/control surfaces and are not sources of trading truth.

## MVP non-goals

The MVP excludes leverage/perpetual hedging, maker strategies, cross-chain bridging as part of execution, flash loans, third-party funds, automated withdrawals/rebalancing, AI-generated live orders, and simultaneous non-atomic leg execution.

## Consequences

This ADR intentionally accepts lower capital efficiency and fewer venues in exchange for deterministic reconciliation and bounded operational risk. Later optimizations must preserve the separation between observation, reservation, external order state, chain state, recovery, and realized accounting.

## Implementation gate

Issue #3 and later implementation issues may proceed using this ADR. A feature is not considered live-ready merely because code exists; venue capability verification, shadow/replay evidence, security tests, explicit Owner limits, and recovery readiness remain separate gates.
