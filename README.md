# idpproxy

idpproxy は、OAuth / OIDC による外部 IdP（GitHub など）の認証を中継・正規化するための  
**認証中継 API（IdP Proxy API）** です。

フロントエンド（Web / SPA / モバイル）と IdP の間に位置し、

- 認証フローの共通化
- ID トークン / セッション管理
- 将来的な IdP 追加・切り替え

を容易にすることを目的としています。

現在は **途中公開（Work In Progress）** の状態であり、  
構成を中心に公開しています。
