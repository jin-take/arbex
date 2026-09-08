# arbex

Arbitrage × Execution

Goの取引エンジンで外部LP・取引所・DeFiの流動性を合成板へ統合するアービトラージシステムです。接続先別の在庫・資金予約と発注ルーター、部分約定・未約定側のヘッジ/回復を設計し、Next.jsの監視UIとSolidityのオンチェーン実行を組み合わせます。

現在は概要設計を議論する段階です。アプリケーション・実売買はまだ実装していません。

## 設計

- [概要設計（議論用ドラフト）](docs/design/overview.md)
- [開発ロードマップ](https://github.com/jin-take/arbex/issues/1)
- [概要設計Issue](https://github.com/jin-take/arbex/issues/2)

設計の変更は `feature/overview-design` ブランチのドラフトPRに集約します。概要設計の「相談する論点（D01〜D12）」を起点に、同じブランチ上で更新します。
