---
id: 1911
title: "ポーカー"
slug: "poker"
status: publish
parent: 0
menu_order: 0
---

<style>
/* 固定ページの投稿日/更新日をこのページのみ非表示 */
.page-id-1911 .entry-meta,
.page-id-1911 .post-meta,
.page-id-1911 .post-info,
.page-id-1911 .post-date,
.page-id-1911 .post-update,
.page-id-1911 .entry-date,
.page-id-1911 .meta-date,
.page-id-1911 .meta-updated,
.page-id-1911 .posted-on,
.page-id-1911 .updated,
.page-id-1911 .date,
.page-id-1911 .modified,
.page-id-1911 time {
  display: none !important;
}
</style>

<!-- CDN: Google Fonts -->
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;700&display=swap" rel="stylesheet">

<!-- CDN: Font Awesome -->
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css">

<!-- CDN: Animate.css -->
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/animate.css/4.1.1/animate.min.css">

<style>
/* ============================================
   CSS Variables - 統一カラー・フォント・スペーシング
   ============================================ */
.poker-timer-app {
  /* Colors */
  --color-primary: #3b82f6;
  --color-primary-hover: #2563eb;
  --color-success: #22c55e;
  --color-success-dark: #16a34a;
  --color-warning: #f59e0b;
  --color-warning-dark: #d97706;
  --color-danger: #ef4444;
  --color-gold: #fbbf24;
  --color-text: #1e293b;
  --color-text-muted: #475569;
  --color-text-light: #64748b;
  --color-text-lighter: #94a3b8;
  --color-bg: #f8fafc;
  --color-card: #fff;
  --color-border: #e2e8f0;
  --color-border-light: #f1f5f9;

  /* Typography - Desktop */
  --font-timer: 96px;
  --font-blind: 42px;
  --font-blind-next: 28px;
  --font-label: 16px;
  --font-value: 32px;
  --font-level-badge: 22px;

  /* Spacing */
  --panel-padding: 20px;
  --card-radius: 12px;
  --card-radius-lg: 16px;
  --gap-sm: 8px;
  --gap-md: 12px;
  --gap-lg: 15px;

  /* Fullscreen Typography Scales */
  --fs-timer-full: 180px;
  --fs-blind-full: 72px;
  --fs-level-badge-full: 28px;
  --fs-label-full: 28px;
  --fs-value-full: 48px;

  /* Mobile Fullscreen Typography Scales */
  --fs-timer-mobile: 100px;
  --fs-blind-mobile: 40px;
  --fs-level-badge-mobile: 20px;
  --fs-label-mobile: 22px;
  --fs-value-mobile: 36px;
}

/* メインアプリ */
.poker-timer-app {
  max-width: 1100px;
  margin: 0 auto;
  font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
  background: var(--color-bg);
  color: var(--color-text);
  border-radius: var(--card-radius-lg);
  padding: var(--panel-padding);
  box-shadow: 0 4px 20px rgba(0,0,0,0.1);
  position: relative;
  touch-action: manipulation;
}

.poker-timer-app * {
  box-sizing: border-box;
}

/* 3カラムレイアウト */
.timer-grid {
  display: grid;
  grid-template-columns: 220px 1fr 240px;
  gap: var(--gap-lg);
  min-height: 400px;
}

/* 左カラム - PRIZE */
.left-panel {
  background: var(--color-card);
  border-radius: var(--card-radius);
  padding: 15px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  display: flex;
  flex-direction: column;
}

.prize-header {
  border-bottom: 2px solid var(--color-border);
  padding-bottom: 16px;
  margin-bottom: 16px;
}

.prize-inmoney {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-muted);
  margin-top: 8px;
}

.prize-inmoney-value {
  font-size: 22px;
  font-weight: 800;
  color: inherit;
}

.prize-list {
  flex: 1;
  overflow: hidden;
  position: relative;
  /* max-height削除: 画像の上まで自動的に伸びる */
}

.prize-list-inner {
  position: absolute;
  width: 100%;
}

.prize-item {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  font-size: 18px;
  border-bottom: 1px solid #f1f5f9;
}

.prize-item:last-child {
  border-bottom: none;
}

.prize-rank {
  color: var(--color-text-muted);
  font-weight: 600;
}

.prize-amount {
  color: var(--color-success);
  font-weight: 700;
}

/* 画像領域（プライズ下に固定表示） */
.mascot-area {
  display: block;
  margin-top: auto;  /* 下に固定 */
  padding-top: 12px;
  text-align: center;
}

.mascot-img {
  width: 180px;
  height: auto;
  border-radius: var(--card-radius);
  opacity: 0.95;
}

/* 中央カラム - タイマー */
.center-panel {
  background: linear-gradient(135deg, var(--color-text) 0%, #334155 100%);
  border-radius: var(--card-radius);
  padding: 25px;
  color: #fff;
  text-align: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
  position: relative;
}

/* レベル表示 - 大きく見やすく */
.level-badge {
  display: inline-block;
  background: var(--color-primary);
  color: #fff;
  padding: 12px 32px;
  border-radius: 30px;
  font-size: var(--font-level-badge);
  font-weight: 700;
  margin-bottom: 15px;
  letter-spacing: 1px;
}

.level-badge.break {
  background: var(--color-success);
}

/* タイマー表示 - 大きく見やすく */
.timer-time {
  font-size: var(--font-timer);
  font-weight: 700;
  font-family: 'JetBrains Mono', 'Courier New', monospace;
  letter-spacing: 2px;
  margin: 20px 0;
  color: #ffffff;
  text-shadow: 0 4px 15px rgba(0,0,0,0.4);
}

.timer-time.warning {
  color: var(--color-danger);
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}

/* プログレスバー */
.progress-bar {
  width: 100%;
  height: 10px;
  background: rgba(255,255,255,0.2);
  border-radius: 5px;
  margin: 20px 0;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #fbbf24, #f59e0b);
  transition: width 1s linear;
  border-radius: 5px;
}

/* ブラインド情報 - 視認性向上 */
.blind-info {
  margin: 15px 0;
}

.blind-current {
  font-size: var(--font-blind);
  font-weight: 700;
  color: var(--color-gold);
  margin: 16px 0;
  text-shadow: 0 3px 12px rgba(251, 191, 36, 0.4);
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 20px;
  flex-wrap: nowrap;
  white-space: nowrap;
}

.blind-current .ante-value {
  font-size: var(--font-blind);
  font-weight: 700;
  color: var(--color-gold);
  text-shadow: 0 3px 12px rgba(251, 191, 36, 0.4);
}

.blind-ante {
  display: none;
}

.blind-next {
  font-size: 18px;
  color: var(--color-text-lighter);
  margin-top: 20px;
  padding-top: 20px;
  border-top: 2px solid rgba(255,255,255,0.2);
}

.blind-next-label {
  font-size: var(--font-label);
  font-weight: 600;
  color: #cbd5e1;
  margin-bottom: 6px;
}

.blind-next-value {
  font-size: var(--font-blind-next);
  font-weight: 600;
  color: #e2e8f0;
  display: block;
  margin-top: 6px;
}

.next-ante {
  font-size: 20px;
  font-weight: 600;
  color: #e2e8f0;
  margin-left: 10px;
}

/* コントロールボタン */
.controls {
  display: flex;
  gap: 10px;
  justify-content: center;
  margin-top: 15px;
}

.btn {
  padding: 12px 24px;
  font-size: 15px;
  font-weight: 600;
  border: none;
  border-radius: var(--gap-sm);
  cursor: pointer;
  transition: all 0.2s;
}

.label-desktop {
  display: inline;
}

.label-mobile {
  display: none;
}

.btn:hover {
  transform: translateY(-2px);
}

.btn-primary {
  background: var(--color-success);
  color: #fff;
}

.btn-primary:hover {
  background: var(--color-success-dark);
}

.btn-secondary {
  background: var(--color-text-light);
  color: #fff;
}

.btn-secondary:hover {
  background: var(--color-text-muted);
}

.btn-warning {
  background: var(--color-warning);
  color: #fff;
}

.btn-warning:hover {
  background: var(--color-warning-dark);
}

/* 右カラム */
.right-panel {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
  position: relative;
  padding-top: 40px;
  height: 100%;
}

/* スペーサー - 中央の空きを埋める */
.right-panel-spacer {
  flex: 1;
}

/* 下部固定: STACK + PLAYERS */
.right-panel-bottom {
  display: flex;
  flex-direction: column;
  gap: var(--gap-md);
}

/* フルスクリーンボタン共通スタイル */
.fullscreen-btn-top,
.fullscreen-btn-mobile {
  color: #fff;
  border: none;
  border-radius: var(--gap-sm);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}

/* フルスクリーンボタン（PC用・right-panel右上固定） */
.fullscreen-btn-top {
  position: absolute;
  top: 0;
  right: 0;
  width: 36px;
  height: 36px;
  background: var(--color-primary);
  font-size: 16px;
  z-index: 10;
}

.fullscreen-btn-top:hover {
  background: var(--color-primary-hover);
}

/* フルスクリーンボタン（モバイル用・center-panelの右上に絶対配置） */
.fullscreen-btn-mobile {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 32px;
  height: 32px;
  background: rgba(59, 130, 246, 0.9);
  font-size: 14px;
  display: none;
  z-index: 20;
}

.fullscreen-btn-mobile:hover {
  background: var(--color-primary-hover);
}

/* モバイル用: Levelバッジ行 */
.level-row {
  display: none;
  align-items: center;
  justify-content: center;
  margin-bottom: 6px;
}

.level-row .level-badge {
  margin-bottom: 0;
  margin-top: 0;
}

.info-card {
  background: var(--color-card);
  border-radius: var(--card-radius-lg);
  padding: var(--panel-padding);
  text-align: center;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
  overflow: hidden;
  min-width: 0;
}

/* ラベルスタイル - 大きく見やすく */
.panel-label {
  font-size: var(--font-label);
  font-weight: 700;
  color: var(--color-text);
  text-transform: uppercase;
  letter-spacing: 1.5px;
  margin-bottom: 10px;
}

.panel-value {
  font-size: var(--font-value);
  font-weight: 700;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.panel-value.gold {
  color: var(--color-warning-dark);
}

.panel-value.large {
  font-size: 40px;
}

/* NEXT BREAK IN */
.break-card {
  background: linear-gradient(135deg, var(--color-success) 0%, var(--color-success-dark) 100%);
  color: #fff;
  padding: 24px;
}

.break-card .panel-label {
  color: rgba(255,255,255,0.9);
}

.break-card .panel-value {
  color: #fff;
  font-size: 42px;
  font-weight: 700;
}

/* ブレイクがない場合のスタイル */
.break-card.no-break {
  background: linear-gradient(135deg, var(--color-text-light) 0%, var(--color-text-muted) 100%);
}

.break-card.no-break .panel-value {
  color: rgba(255,255,255,0.6);
}

/* STACK カード */
.stack-card {
  padding: var(--panel-padding);
}

.stack-card .stack-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--gap-sm) 0;
}

.stack-card .stack-label {
  font-size: var(--font-label);
  font-weight: 500;
  color: var(--color-text-muted);
}

.stack-card .stack-value {
  font-size: var(--font-value);
  font-weight: 700;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* PLAYERS カード */
.players-card {
  padding: var(--panel-padding);
}

.players-card .players-display {
  font-size: var(--font-value);
  font-weight: 700;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin: var(--gap-md) 0;
}

.players-card .players-label {
  font-size: 14px;
  color: var(--color-text-muted);
  margin-bottom: var(--gap-md);
}

/* カウンターコントロール */
.counter-controls {
  display: flex;
  justify-content: center;
  gap: var(--gap-sm);
}

.counter-btn {
  width: 36px;
  height: 36px;
  border-radius: var(--gap-sm);
  border: 2px solid var(--color-border);
  background: var(--color-card);
  color: var(--color-text-muted);
  font-size: var(--font-label);
  font-weight: bold;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.counter-btn:hover {
  background: var(--color-primary);
  border-color: var(--color-primary);
  color: #fff;
}

/* 設定モーダル - z-index最大 */
.settings-modal {
  display: none;
  position: fixed !important;
  top: 0 !important;
  left: 0 !important;
  width: 100vw !important;
  height: 100vh !important;
  background: rgba(0,0,0,0.95) !important;
  z-index: 2147483647 !important;
  overflow-y: auto;
  isolation: isolate;
}

.settings-modal.active {
  display: flex !important;
  justify-content: center;
  align-items: center;
  padding: 20px;
  box-sizing: border-box;
}

/* モーダル表示時にWordPressのヘッダー/フッターを非表示 */
body.modal-open header,
body.modal-open .site-header,
body.modal-open #masthead,
body.modal-open .header,
body.modal-open footer,
body.modal-open .site-footer,
body.modal-open #colophon,
body.modal-open .footer {
  display: none !important;
}

.settings-content {
  background: #ffffff !important;
  border-radius: 16px;
  padding: 24px;
  max-width: 650px;
  width: 100%;
  margin: 0;
  color: #1e293b !important;
  max-height: calc(100vh - 40px);
  overflow-y: auto;
  position: relative;
  z-index: 1;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

.settings-title,
.settings-modal h2.settings-title,
.settings-content h2.settings-title {
  font-size: 20px !important;
  font-weight: bold !important;
  margin-bottom: 20px !important;
  text-align: center !important;
  color: #1e293b !important;
  background: transparent !important;
  padding: 0 !important;
  border: none !important;
  margin-top: 0 !important;
}

.settings-tabs {
  display: flex;
  border-bottom: 2px solid #e2e8f0;
  margin-bottom: 20px;
}

.settings-tab {
  flex: 1;
  padding: 10px;
  text-align: center;
  cursor: pointer;
  font-weight: 600;
  color: #94a3b8 !important;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
  background: transparent !important;
}

.settings-tab.active {
  color: #3b82f6 !important;
  border-bottom-color: #3b82f6;
}

.settings-panel {
  display: none;
}

.settings-panel.active {
  display: block;
}

.setting-group {
  margin-bottom: 18px;
}

.setting-label {
  display: block;
  margin-bottom: 6px;
  font-weight: 600;
  color: #475569 !important;
  font-size: 14px;
}

.setting-input {
  width: 100%;
  padding: 10px;
  border: 2px solid #e2e8f0 !important;
  border-radius: 8px;
  background: #f8fafc !important;
  color: #1e293b !important;
  font-size: 16px;
  box-sizing: border-box;
}

.setting-input:focus {
  outline: none;
  border-color: #3b82f6 !important;
}

.setting-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

/* ブラインドセット選択 */
.blind-set-selector {
  display: flex;
  gap: 8px;
  margin-bottom: 15px;
}

.blind-set-selector select {
  flex: 1;
  padding: 10px;
  border: 2px solid #e2e8f0 !important;
  border-radius: 8px;
  font-size: 16px;
  background: #f8fafc !important;
  color: #1e293b !important;
}

.blind-set-selector button {
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
}

.btn-save-set {
  background: #22c55e !important;
  color: #fff !important;
}

.btn-delete-set {
  background: #ef4444 !important;
  color: #fff !important;
}

/* ブラインドレベル一覧 */
.blind-levels {
  max-height: 280px;
  overflow-y: auto;
  margin-bottom: 16px;
  border: 1px solid #e2e8f0 !important;
  border-radius: 8px;
}

/* 共通ドラッグアイテムスタイル */
.blind-level-item,
.prize-edit-item {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 8px;
  background: #f8fafc !important;
  border-bottom: 1px solid #e2e8f0 !important;
  cursor: grab;
  transition: background 0.2s, transform 0.2s;
  flex-wrap: wrap;
}

.blind-level-item:active,
.prize-edit-item:active {
  cursor: grabbing;
}

.blind-level-item.dragging,
.prize-edit-item.dragging {
  opacity: 0.5;
  background: #e2e8f0 !important;
}

.blind-level-item.drag-over,
.prize-edit-item.drag-over {
  border-top: 3px solid #3b82f6 !important;
}

.blind-level-item:last-child,
.prize-edit-item:last-child {
  border-bottom: none !important;
}

.blind-level-item.break-item {
  background: #dcfce7 !important;
}

.blind-level-item input {
  width: 84px;
  padding: 4px;
  border: 1px solid #e2e8f0 !important;
  border-radius: 4px;
  background: #ffffff !important;
  color: #1e293b !important;
  font-size: 16px;
  box-sizing: border-box;
  text-align: center;
}

.blind-level-item input.time-input {
  width: 56px;
}

.blind-level-item input.break-time-input {
  width: 56px !important;
}

.blind-level-item input.sb-input,
.blind-level-item input.bb-input,
.blind-level-item input.ante-input {
  width: 84px;
}

/* 共通ドラッグハンドル */
.blind-level-item .drag-handle,
.prize-edit-item .drag-handle {
  cursor: grab;
  color: #94a3b8 !important;
  font-size: 14px;
  padding: 0 4px;
}

.blind-level-item .drag-handle:active,
.prize-edit-item .drag-handle:active {
  cursor: grabbing;
}

.blind-level-item .level-num {
  min-width: 24px;
  font-weight: bold;
  color: #3b82f6 !important;
  font-size: 13px;
  text-align: center;
}

.blind-level-item .break-label {
  color: #22c55e !important;
  font-weight: 600;
  font-size: 14px;
}

/* 共通削除ボタン */
.blind-level-item .delete-level,
.prize-edit-item .delete-prize {
  background: #ef4444 !important;
  color: #fff !important;
  border: none;
  border-radius: 4px;
  padding: 4px 8px;
  cursor: pointer;
  font-size: 14px;
  margin-left: auto;
}

.move-controls {
  display: none;
  gap: 4px;
  flex-direction: column;
}

.move-btn {
  background: #e2e8f0 !important;
  color: #1e293b !important;
  border: none;
  border-radius: 6px;
  width: 28px;
  height: 28px;
  font-size: 14px;
  cursor: pointer;
}

.move-btn:active {
  transform: translateY(1px);
}

.blind-level-item .time-label {
  font-size: 12px;
  color: #94a3b8 !important;
}

.blind-level-item span {
  font-size: 12px;
  color: #64748b !important;
}

.level-actions {
  display: flex;
  gap: 8px;
}

.level-actions button {
  flex: 1;
  padding: 10px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
}

.btn-add-level {
  background: #3b82f6 !important;
  color: #fff !important;
}

.btn-add-break {
  background: #22c55e !important;
  color: #fff !important;
}

/* プライズ編集リスト */
.prize-edit-list {
  max-height: 250px;
  overflow-y: auto;
  border: 1px solid #e2e8f0 !important;
  border-radius: 8px;
  margin-top: 10px;
}

/* prize-edit-item の基本スタイルは共通ドラッグアイテムスタイルで定義済み */

.prize-edit-item input {
  padding: 4px;
  border: 1px solid #e2e8f0 !important;
  border-radius: 4px;
  font-size: 16px;
  box-sizing: border-box;
  text-align: center;
  background: #ffffff !important;
  color: #1e293b !important;
}

.prize-edit-item input.rank-input {
  width: 56px;
}

.prize-edit-item input.amount-input {
  width: 104px;
  text-align: right;
}

.prize-edit-item span {
  font-size: 13px;
  color: #64748b !important;
}

/* delete-prize は共通削除ボタンで定義済み */

.prize-edit-actions {
  display: flex;
  gap: 8px;
  margin-top: 10px;
}

.btn-add-prize {
  flex: 1;
  padding: 10px;
  background: #3b82f6 !important;
  color: #fff !important;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
}

.prize-calc-row {
  display: flex;
  gap: 10px;
  align-items: center;
  margin-bottom: 10px;
}

.btn-calc {
  background: #3b82f6 !important;
  color: #fff !important;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
}

/* モーダルボタン */
.modal-buttons {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.modal-buttons .btn {
  flex: 1;
}

/* モーダル内ボタンの色を明示的に指定 */
.settings-modal .btn-primary {
  background: #22c55e !important;
  color: #fff !important;
}

.settings-modal .btn-primary:hover {
  background: #16a34a !important;
}

.settings-modal .btn-secondary {
  background: #64748b !important;
  color: #fff !important;
}

.settings-modal .btn-secondary:hover {
  background: #475569 !important;
}

/* ============================================
   共通フルスクリーンスタイル
   :fullscreen と .mobile-fullscreen で共有
   ============================================ */

/* フルスクリーン共通: 基本レイアウト */
.poker-timer-app:fullscreen,
.poker-timer-app.mobile-fullscreen {
  display: flex;
  flex-direction: column;
  max-width: none !important;
  border-radius: 0 !important;
  background: var(--color-bg) !important;
  box-sizing: border-box !important;
  overflow: hidden !important;
}

/* フルスクリーン共通: 設定ボタン非表示 */
.poker-timer-app:fullscreen .controls #btnSettings,
.poker-timer-app.mobile-fullscreen .controls #btnSettings {
  display: none !important;
}

/* フルスクリーン共通: フルスクリーンボタン非表示 */
.poker-timer-app:fullscreen .fullscreen-btn-top,
.poker-timer-app.mobile-fullscreen .fullscreen-btn-top,
.poker-timer-app.mobile-fullscreen .fullscreen-btn-mobile {
  display: none;
}

/* フルスクリーン共通: timer-grid */
.poker-timer-app:fullscreen .timer-grid,
.poker-timer-app.mobile-fullscreen .timer-grid {
  flex: 1;
  width: 100%;
  min-height: 0;
  box-sizing: border-box;
}

/* フルスクリーン共通: パネル */
.poker-timer-app:fullscreen .left-panel,
.poker-timer-app.mobile-fullscreen .left-panel,
.poker-timer-app:fullscreen .right-panel,
.poker-timer-app.mobile-fullscreen .right-panel {
  min-width: 0;
  overflow: hidden;
}

/* フルスクリーン共通: 中央パネル */
.poker-timer-app:fullscreen .center-panel,
.poker-timer-app.mobile-fullscreen .center-panel {
  min-width: 0;
  max-width: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  overflow: hidden;
}

/* フルスクリーン共通: プレイヤーボタン非表示 */
.poker-timer-app:fullscreen .players-buttons,
.poker-timer-app.mobile-fullscreen .players-buttons {
  display: none;
}

/* フルスクリーン共通: テキストオーバーフロー対策 */
.poker-timer-app:fullscreen .panel-value,
.poker-timer-app:fullscreen .stack-value,
.poker-timer-app:fullscreen .players-display,
.poker-timer-app.mobile-fullscreen .panel-value,
.poker-timer-app.mobile-fullscreen .stack-value,
.poker-timer-app.mobile-fullscreen .players-display {
  white-space: nowrap !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
}

/* ============================================
   デスクトップ ネイティブフルスクリーン固有スタイル
   ============================================ */
.poker-timer-app:fullscreen {
  padding: 40px;
  width: 100vw !important;
  height: 100vh !important;
}

.poker-timer-app:fullscreen .timer-grid {
  height: 100%;
  gap: 16px;
  grid-template-columns: minmax(180px, 1fr) 2.5fr minmax(180px, 1fr);
  padding: 16px;
}

/* デスクトップフルスクリーン: タイマー */
.poker-timer-app:fullscreen .timer-time {
  font-size: var(--fs-timer-full);
  font-weight: 700;
  color: #ffffff;
  text-shadow: 0 4px 20px rgba(0,0,0,0.3);
}

.poker-timer-app:fullscreen .timer-time.warning,
.poker-timer-app.mobile-fullscreen .timer-time.warning {
  color: #ef4444 !important;
  animation: pulse 1s infinite;
}

.poker-timer-app:fullscreen .timer-level {
  font-size: 36px;
  font-weight: 600;
  color: #ffffff;
}

.poker-timer-app:fullscreen .level-badge {
  font-size: var(--fs-level-badge-full);
  padding: 14px 36px;
}

.poker-timer-app:fullscreen .blind-current {
  font-size: var(--fs-blind-full);
  font-weight: 700;
  color: var(--color-gold);
}

.poker-timer-app:fullscreen .blind-current .ante-value {
  font-size: var(--fs-blind-full);
  color: var(--color-gold);
}

.poker-timer-app:fullscreen .blind-next {
  margin-top: 16px;
}

.poker-timer-app:fullscreen .blind-next-label {
  font-size: 22px;
  color: var(--color-text-lighter);
}

.poker-timer-app:fullscreen .blind-next-value {
  font-size: 42px;
  font-weight: 600;
  color: #e2e8f0;
}

.poker-timer-app:fullscreen .next-ante {
  font-size: 28px;
  color: #e2e8f0;
}

.poker-timer-app:fullscreen .timer-progress {
  height: 16px;
  margin: 24px 0;
}

/* デスクトップフルスクリーン: 左パネル - PRIZE */
.poker-timer-app:fullscreen .left-panel {
  height: 100%;
  padding: 24px;
}

.poker-timer-app:fullscreen .left-panel .panel-label {
  font-size: var(--fs-label-full);
  font-weight: 600;
  color: var(--color-text);
}

.poker-timer-app:fullscreen .prize-inmoney {
  font-size: 26px;
}

.poker-timer-app:fullscreen .prize-inmoney-value {
  font-size: 26px;
}

.poker-timer-app:fullscreen .prize-item {
  font-size: 24px;
  padding: 6px 0;
}

.poker-timer-app:fullscreen .prize-item .prize-rank {
  font-size: 24px;
  color: var(--color-text);
}

.poker-timer-app:fullscreen .prize-item .prize-amount {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-success);
}

/* デスクトップフルスクリーン: 右パネル */
.poker-timer-app:fullscreen .break-card {
  padding: var(--panel-padding);
}

.poker-timer-app:fullscreen .break-card .panel-label {
  font-size: var(--fs-label-full);
  font-weight: 600;
}

.poker-timer-app:fullscreen .break-card .panel-value {
  font-size: 36px;
  font-weight: 700;
}

.poker-timer-app:fullscreen .stack-card {
  padding: var(--panel-padding);
}

.poker-timer-app:fullscreen .stack-card .panel-label {
  font-size: var(--fs-label-full);
  font-weight: 600;
  color: var(--color-text);
}

.poker-timer-app:fullscreen .stack-row .stack-label {
  font-size: 24px;
  color: var(--color-text-muted);
}

.poker-timer-app:fullscreen .stack-row .stack-value {
  font-size: var(--fs-value-full);
  font-weight: 700;
  color: var(--color-text);
}

.poker-timer-app:fullscreen .players-card {
  padding: var(--panel-padding);
}

.poker-timer-app:fullscreen .players-card .panel-label {
  font-size: var(--fs-label-full);
  font-weight: 600;
  color: var(--color-text);
}

.poker-timer-app:fullscreen .players-display {
  font-size: var(--fs-value-full);
  font-weight: 700;
  color: var(--color-text);
}

.poker-timer-app:fullscreen .players-label {
  font-size: 22px;
  color: var(--color-text-muted);
}

.poker-timer-app:fullscreen .mascot-img {
  width: 280px;
  height: auto;
}

/* ============================================
   モバイル擬似フルスクリーン固有スタイル
   ============================================ */
body.mobile-fullscreen-active {
  overflow: hidden !important;
}

/* フルスクリーン時にWordPressのヘッダー/フッターを非表示 */
body.mobile-fullscreen-active header,
body.mobile-fullscreen-active #masthead,
body.mobile-fullscreen-active .site-header,
body.mobile-fullscreen-active .header,
body.mobile-fullscreen-active nav,
body.mobile-fullscreen-active .navigation,
body.mobile-fullscreen-active .main-navigation,
body.mobile-fullscreen-active footer,
body.mobile-fullscreen-active #colophon,
body.mobile-fullscreen-active .site-footer,
body.mobile-fullscreen-active .footer,
body.mobile-fullscreen-active .sidebar,
body.mobile-fullscreen-active aside {
  display: none !important;
  visibility: hidden !important;
}

.poker-timer-app.mobile-fullscreen {
  position: fixed !important;
  top: 0 !important;
  left: 0 !important;
  right: 0 !important;
  width: auto !important;
  max-width: 100% !important;
  z-index: 2147483647 !important;
  padding: 8px;
  padding-top: max(8px, env(safe-area-inset-top));
  padding-bottom: max(8px, env(safe-area-inset-bottom));
  padding-left: max(8px, env(safe-area-inset-left));
  padding-right: max(8px, env(safe-area-inset-right));
}

.poker-timer-app.mobile-fullscreen .timer-grid {
  max-width: 100%;
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: auto 1fr auto;
  padding: 4px;
  gap: 4px;
  align-items: stretch;
}

/* ============================================
   iPad以上（横幅901px以上）3カラムレイアウト
   ============================================ */
@media (min-width: 901px) {
  /* Flexboxで3カラムを実現（iPad Safari対応） */
  .poker-timer-app.mobile-fullscreen .timer-grid {
    display: flex !important;
    flex-direction: row !important;
    flex-wrap: nowrap !important;
    gap: 0 !important;
    padding: 8px !important;
    flex: 1 1 auto !important;
    height: calc(100% - 16px) !important;
    max-height: calc(100% - 16px) !important;
  }

  /* 3カラム共通幅設定（marginでgap代替） - 中央パネルを広く */
  .poker-timer-app.mobile-fullscreen .left-panel {
    flex: 0 0 calc(20% - 6px) !important;
    width: calc(20% - 6px) !important;
    max-width: calc(20% - 6px) !important;
    margin-right: 8px !important;
  }

  .poker-timer-app.mobile-fullscreen .center-panel {
    flex: 0 0 calc(60% - 4px) !important;
    width: calc(60% - 4px) !important;
    max-width: calc(60% - 4px) !important;
    margin-right: 8px !important;
  }

  .poker-timer-app.mobile-fullscreen .right-panel {
    flex: 0 0 calc(20% - 6px) !important;
    width: calc(20% - 6px) !important;
    max-width: calc(20% - 6px) !important;
    margin-right: 0 !important;
    display: flex !important;
    flex-direction: column !important;
    justify-content: flex-end !important;
    gap: var(--gap-md) !important;
  }

  /* オーバーフロー制御 */
  .poker-timer-app.mobile-fullscreen .timer-grid *,
  .poker-timer-app.mobile-fullscreen .timer-grid *::before,
  .poker-timer-app.mobile-fullscreen .timer-grid *::after {
    min-width: 0 !important;
    max-width: 100% !important;
    box-sizing: border-box !important;
  }

  .poker-timer-app.mobile-fullscreen .info-card {
    min-width: 0 !important;
    max-width: 100% !important;
    overflow: hidden !important;
  }

  .poker-timer-app.mobile-fullscreen .right-panel .info-card {
    margin-bottom: 10px !important;
  }
  .poker-timer-app.mobile-fullscreen .right-panel .info-card:last-child {
    margin-bottom: 0 !important;
  }

  /* iPad: prize-list表示 */
  .poker-timer-app.mobile-fullscreen .prize-list {
    flex: 1 !important;
  }

  .poker-timer-app.mobile-fullscreen .prize-list-inner {
    position: relative !important;
  }

  .poker-timer-app.mobile-fullscreen .prize-item {
    display: flex !important;
    justify-content: space-between !important;
    font-size: 20px !important;
    padding: 4px 0 !important;
    border-bottom: 1px solid rgba(0,0,0,0.1) !important;
  }

  /* iPad: フォントサイズ */
  .poker-timer-app.mobile-fullscreen .prize-header .panel-label { font-size: 20px !important; }
  .poker-timer-app.mobile-fullscreen .prize-inmoney { font-size: 16px !important; }
  .poker-timer-app.mobile-fullscreen .prize-inmoney-value { font-size: 20px !important; }
  .poker-timer-app.mobile-fullscreen .timer-time { font-size: 130px !important; line-height: 1 !important; }
  .poker-timer-app.mobile-fullscreen .timer-level { font-size: 28px !important; }
  .poker-timer-app.mobile-fullscreen .level-badge { font-size: 24px !important; padding: 8px 20px !important; }
  .poker-timer-app.mobile-fullscreen .blind-current { font-size: 52px !important; }
  .poker-timer-app.mobile-fullscreen .blind-current .ante-value { font-size: 52px !important; }
  .poker-timer-app.mobile-fullscreen .blind-next-label { font-size: 20px !important; }
  .poker-timer-app.mobile-fullscreen .blind-next-value { font-size: 32px !important; }
  .poker-timer-app.mobile-fullscreen .next-ante { font-size: 20px !important; }
  .poker-timer-app.mobile-fullscreen .timer-progress { height: 10px !important; margin: 8px 0 !important; }

  /* iPad: QRコード画像 */
  .poker-timer-app.mobile-fullscreen .mascot-area {
    display: block !important;
    margin-top: auto !important;
    padding-top: 8px !important;
  }

  .poker-timer-app.mobile-fullscreen .mascot-img {
    width: 240px !important;
    height: auto !important;
  }
}

/* モバイルフルスクリーン: デフォルトフォントサイズ */
.poker-timer-app.mobile-fullscreen .timer-card {
  padding: 6px;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 0;
}

.poker-timer-app.mobile-fullscreen .timer-time {
  font-size: var(--fs-timer-mobile);
  font-weight: 700;
  color: #ffffff;
  text-shadow: 0 4px 20px rgba(0,0,0,0.3);
  line-height: 1;
}

.poker-timer-app.mobile-fullscreen .timer-level {
  font-size: 22px;
  font-weight: 600;
  color: #ffffff;
}

.poker-timer-app.mobile-fullscreen .level-badge {
  padding: 6px 16px;
  font-size: var(--fs-level-badge-mobile);
}

.poker-timer-app.mobile-fullscreen .blind-current {
  font-size: var(--fs-blind-mobile);
  font-weight: 700;
  color: var(--color-gold);
}

.poker-timer-app.mobile-fullscreen .blind-current .ante-value {
  font-size: var(--fs-blind-mobile);
  color: var(--color-gold);
}

.poker-timer-app.mobile-fullscreen .blind-next-label { font-size: 16px; }
.poker-timer-app.mobile-fullscreen .blind-next-value { font-size: 24px; font-weight: 600; color: #e2e8f0; }
.poker-timer-app.mobile-fullscreen .next-ante { font-size: 16px; color: #e2e8f0; }
.poker-timer-app.mobile-fullscreen .timer-progress { height: 6px; margin: 4px 0; }

/* モバイルフルスクリーン: 左パネル */
.poker-timer-app.mobile-fullscreen .left-panel {
  height: 100%;
}

.poker-timer-app.mobile-fullscreen .left-panel .panel-label { font-size: var(--fs-label-mobile); margin-bottom: 2px; }
.poker-timer-app.mobile-fullscreen .prize-inmoney { font-size: var(--fs-label-mobile); }
.poker-timer-app.mobile-fullscreen .prize-inmoney-value { font-size: 28px; }
.poker-timer-app.mobile-fullscreen .prize-item { font-size: 24px; padding: 3px 0; }

/* モバイルフルスクリーン: 右パネル */
.poker-timer-app.mobile-fullscreen .right-panel {
  gap: 2px;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
}

.poker-timer-app.mobile-fullscreen .right-panel-spacer { display: none !important; }

.poker-timer-app.mobile-fullscreen .right-panel-bottom {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 0 0 auto;
  min-height: 0;
}

/* モバイルフルスクリーン: カード共通 */
.poker-timer-app.mobile-fullscreen .break-card,
.poker-timer-app.mobile-fullscreen .stack-card,
.poker-timer-app.mobile-fullscreen .players-card {
  padding: 4px 6px;
  flex: 0 1 auto;
  min-height: 0;
}

.poker-timer-app.mobile-fullscreen .break-card .panel-label,
.poker-timer-app.mobile-fullscreen .stack-card .panel-label,
.poker-timer-app.mobile-fullscreen .players-card .panel-label {
  font-size: var(--fs-label-mobile);
  margin-bottom: 2px;
}

.poker-timer-app.mobile-fullscreen .break-card .panel-value { font-size: var(--fs-value-mobile); }
.poker-timer-app.mobile-fullscreen .stack-row .stack-label { font-size: var(--fs-label-mobile); }
.poker-timer-app.mobile-fullscreen .stack-row .stack-value { font-size: var(--fs-value-mobile); }
.poker-timer-app.mobile-fullscreen .stack-row { margin: 1px 0; }
.poker-timer-app.mobile-fullscreen .players-display { font-size: 44px; font-weight: 700; }
.poker-timer-app.mobile-fullscreen .players-label { font-size: 20px; color: var(--color-text-muted); }

/* フルスクリーン解除ボタン */
.exit-fullscreen-btn {
  display: none;
  position: fixed;
  top: max(8px, env(safe-area-inset-top));
  right: max(8px, env(safe-area-inset-right));
  width: 40px;
  height: 40px;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  border: none;
  border-radius: var(--gap-sm);
  font-size: 18px;
  cursor: pointer;
  z-index: 2147483647;
  align-items: center;
  justify-content: center;
}

.poker-timer-app.mobile-fullscreen .exit-fullscreen-btn {
  display: flex;
}

.exit-fullscreen-btn:hover {
  background: rgba(0, 0, 0, 0.8);
}

/* ============================================
   モバイルフルスクリーン: 縦画面（max-width: 599px）
   ============================================ */
@media (max-width: 599px) and (orientation: portrait) {
  .poker-timer-app.mobile-fullscreen .left-panel { display: none !important; }

  .poker-timer-app.mobile-fullscreen .center-panel {
    order: 1;
    align-items: center;
    flex: 1;
  }

  .poker-timer-app.mobile-fullscreen .right-panel {
    order: 2;
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: auto auto;
    gap: 6px;
    padding: 4px;
  }

  .poker-timer-app.mobile-fullscreen .right-panel-bottom { display: contents; }

  .poker-timer-app.mobile-fullscreen .break-card { grid-column: 1; grid-row: 1; padding: 6px 10px; }
  .poker-timer-app.mobile-fullscreen .stack-card { grid-column: 1; grid-row: 2; padding: 6px 10px; }
  .poker-timer-app.mobile-fullscreen .players-card {
    grid-column: 2;
    grid-row: 1 / 3;
    padding: 6px 10px;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  /* フォントサイズ調整 */
  .poker-timer-app.mobile-fullscreen .timer-time { font-size: 64px; }
  .poker-timer-app.mobile-fullscreen .blind-current { font-size: 26px; }
  .poker-timer-app.mobile-fullscreen .blind-current .ante-value { font-size: 26px; }
  .poker-timer-app.mobile-fullscreen .blind-next-value { font-size: 16px; }
  .poker-timer-app.mobile-fullscreen .next-ante { font-size: 16px; }
  .poker-timer-app.mobile-fullscreen .level-badge { font-size: 14px; padding: 3px 10px; }

  .poker-timer-app.mobile-fullscreen .left-panel .panel-label,
  .poker-timer-app.mobile-fullscreen .break-card .panel-label,
  .poker-timer-app.mobile-fullscreen .stack-card .panel-label,
  .poker-timer-app.mobile-fullscreen .players-card .panel-label { font-size: 12px; }

  .poker-timer-app.mobile-fullscreen .break-card .panel-value { font-size: 20px; }
  .poker-timer-app.mobile-fullscreen .stack-row .stack-label { font-size: 12px; }
  .poker-timer-app.mobile-fullscreen .stack-row .stack-value { font-size: 18px; }
  .poker-timer-app.mobile-fullscreen .players-display { font-size: 24px; }
  .poker-timer-app.mobile-fullscreen .players-label { font-size: 11px; }
}

/* ============================================
   モバイルフルスクリーン: 横画面（max-width: 900px）
   ============================================ */
@media (max-width: 900px) and (orientation: landscape) {
  .poker-timer-app.mobile-fullscreen {
    padding: 4px;
    padding-left: max(4px, env(safe-area-inset-left, 4px));
    padding-right: max(4px, env(safe-area-inset-right, 4px));
    padding-bottom: max(4px, env(safe-area-inset-bottom, 4px));
  }

  .poker-timer-app.mobile-fullscreen .timer-grid {
    grid-template-columns: minmax(100px, 1fr) 2fr minmax(100px, 1fr);
    grid-template-rows: 1fr;
    gap: 4px;
    height: 100%;
    max-height: 100%;
  }

  .poker-timer-app.mobile-fullscreen .left-panel {
    display: flex !important;
    flex-direction: column;
    order: 1;
    max-height: 100%;
    padding: 4px;
  }

  .poker-timer-app.mobile-fullscreen .prize-card { padding: 2px 4px; }
  .poker-timer-app.mobile-fullscreen .prize-header { padding-bottom: 2px; margin-bottom: 2px; }
  .poker-timer-app.mobile-fullscreen .prize-inmoney { font-size: 10px; }
  .poker-timer-app.mobile-fullscreen .prize-inmoney-value { font-size: 12px; }
  .poker-timer-app.mobile-fullscreen .prize-list { flex: 1; }
  .poker-timer-app.mobile-fullscreen .prize-list-inner { position: relative; will-change: transform; }
  .poker-timer-app.mobile-fullscreen .prize-item { font-size: 11px; padding: 1px 0; }

  .poker-timer-app.mobile-fullscreen .mascot-area { display: block; margin-top: auto; padding-top: 4px; }
  .poker-timer-app.mobile-fullscreen .mascot-img { width: 180px; height: auto; }

  .poker-timer-app.mobile-fullscreen .center-panel {
    order: 2;
    align-items: center;
    max-height: 100%;
    padding: 4px;
  }

  .poker-timer-app.mobile-fullscreen .right-panel {
    order: 3;
    gap: 2px;
    padding: 4px;
    justify-content: flex-end;
    max-height: 100%;
  }

  .poker-timer-app.mobile-fullscreen .right-panel-bottom { display: contents; }

  .poker-timer-app.mobile-fullscreen .break-card,
  .poker-timer-app.mobile-fullscreen .stack-card,
  .poker-timer-app.mobile-fullscreen .players-card { padding: 2px 4px; }

  /* フォントサイズ調整 */
  .poker-timer-app.mobile-fullscreen .timer-time { font-size: 56px; margin: 2px 0; }
  .poker-timer-app.mobile-fullscreen .blind-current { font-size: 22px; margin: 2px 0; }
  .poker-timer-app.mobile-fullscreen .blind-current .ante-value { font-size: 22px; }
  .poker-timer-app.mobile-fullscreen .blind-next-value { font-size: 14px; }
  .poker-timer-app.mobile-fullscreen .next-ante { font-size: 14px; }
  .poker-timer-app.mobile-fullscreen .level-badge { font-size: 12px; padding: 2px 8px; margin-bottom: 2px; }
  .poker-timer-app.mobile-fullscreen .timer-progress { height: 3px; margin: 2px 0; }
  .poker-timer-app.mobile-fullscreen .blind-info { margin: 2px 0; }
  .poker-timer-app.mobile-fullscreen .blind-next { margin-top: 2px; padding-top: 2px; }

  .poker-timer-app.mobile-fullscreen .left-panel .panel-label,
  .poker-timer-app.mobile-fullscreen .break-card .panel-label,
  .poker-timer-app.mobile-fullscreen .stack-card .panel-label,
  .poker-timer-app.mobile-fullscreen .players-card .panel-label { font-size: 10px; margin-bottom: 1px; }

  .poker-timer-app.mobile-fullscreen .break-card .panel-value { font-size: 16px; }
  .poker-timer-app.mobile-fullscreen .stack-row .stack-label { font-size: 10px; }
  .poker-timer-app.mobile-fullscreen .stack-row .stack-value { font-size: 14px; }
  .poker-timer-app.mobile-fullscreen .stack-row { padding: 1px 0; }
  .poker-timer-app.mobile-fullscreen .players-display { font-size: 20px; }
  .poker-timer-app.mobile-fullscreen .players-label { font-size: 9px; }
}

/* ============================================
   iPad横画面（901px〜1024px）追加調整
   ============================================ */
@media (min-width: 901px) and (max-width: 1024px) and (orientation: landscape) {
  .poker-timer-app.mobile-fullscreen .timer-time { font-size: 80px; }
  .poker-timer-app.mobile-fullscreen .level-badge { font-size: 16px; padding: 6px 14px; margin-bottom: 8px; }
  .poker-timer-app.mobile-fullscreen .blind-current { font-size: 32px; }
  .poker-timer-app.mobile-fullscreen .blind-current .ante-value { font-size: 32px; }
  .poker-timer-app.mobile-fullscreen .blind-next-value { font-size: 18px; }
  .poker-timer-app.mobile-fullscreen .next-ante { font-size: 14px; }
  .poker-timer-app.mobile-fullscreen .timer-progress { height: 6px; margin: 6px 0; }

  .poker-timer-app.mobile-fullscreen .left-panel .panel-label { font-size: 14px; }
  .poker-timer-app.mobile-fullscreen .prize-inmoney { font-size: 14px; }
  .poker-timer-app.mobile-fullscreen .prize-inmoney-value { font-size: 18px; }
  .poker-timer-app.mobile-fullscreen .prize-item { font-size: 16px; padding: 2px 0; }

  .poker-timer-app.mobile-fullscreen .right-panel { gap: 6px; }
  .poker-timer-app.mobile-fullscreen .break-card,
  .poker-timer-app.mobile-fullscreen .stack-card,
  .poker-timer-app.mobile-fullscreen .players-card { padding: 6px 10px; }

  .poker-timer-app.mobile-fullscreen .break-card .panel-label,
  .poker-timer-app.mobile-fullscreen .stack-card .panel-label,
  .poker-timer-app.mobile-fullscreen .players-card .panel-label { font-size: 14px; }

  .poker-timer-app.mobile-fullscreen .break-card .panel-value { font-size: 24px; }
  .poker-timer-app.mobile-fullscreen .stack-row .stack-label { font-size: 14px; }
  .poker-timer-app.mobile-fullscreen .stack-row .stack-value { font-size: 22px; }
  .poker-timer-app.mobile-fullscreen .players-display { font-size: 28px; }
  .poker-timer-app.mobile-fullscreen .players-label { font-size: 13px; }
}

/* 使い方説明 */
.usage-guide {
  max-width: 1100px;
  margin: 30px auto 0;
  padding: 24px;
  background: var(--color-card);
  border-radius: var(--card-radius);
  box-shadow: 0 2px 10px rgba(0,0,0,0.08);
  color: var(--color-text);
  font-size: 14px;
  line-height: 1.8;
}

.usage-guide h3 {
  margin: 0 0 16px 0;
  font-size: 18px;
  color: var(--color-text);
  border-bottom: 2px solid var(--color-border);
  padding-bottom: var(--gap-sm);
}

.usage-guide h4 {
  margin: var(--panel-padding) 0 var(--gap-sm) 0;
  font-size: 15px;
  color: var(--color-text-muted);
}

.usage-guide ul {
  margin: var(--gap-sm) 0;
  padding-left: var(--panel-padding);
}

.usage-guide li {
  margin: 4px 0;
}

.usage-guide kbd {
  display: inline-block;
  padding: 2px 6px;
  font-size: 12px;
  font-family: monospace;
  background: var(--color-border-light);
  border: 1px solid var(--color-border);
  border-radius: 4px;
}

/* ============================================
   レスポンシブ: 縦画面（max-width: 900px）
   ============================================ */
@media (max-width: 900px) and (orientation: portrait) {
  .timer-grid {
    grid-template-columns: 1fr;
    grid-template-rows: auto auto auto;
    gap: var(--gap-md);
  }

  .center-panel {
    order: 1;
    padding: var(--gap-md);
  }

  .level-row { display: flex; margin-bottom: 4px; }
  .center-panel > .level-badge { display: none; }
  .timer-time { margin: 6px 0; font-size: 64px; }
  .blind-info { margin-top: 6px; }
  .controls { margin-top: 10px; }
  .fullscreen-btn-top { display: none; }
  .fullscreen-btn-mobile { display: flex; }

  .right-panel {
    order: 2;
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: auto auto;
    gap: var(--gap-sm);
    padding-top: 0;
    height: auto;
  }

  .right-panel-spacer { display: none; }
  .right-panel-bottom { display: contents; }

  .break-card { grid-column: 1; grid-row: 1; }
  .stack-card { grid-column: 1; grid-row: 2; }
  .players-card {
    grid-column: 2;
    grid-row: 1 / 3;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .left-panel {
    order: 3;
    height: 210px;
    max-height: 210px;
  }

  .level-badge { font-size: 14px; padding: 6px 12px; margin-top: 0; margin-bottom: 0; }
  .blind-current { font-size: 32px; }
  .blind-current .ante-value { font-size: 32px; }
  .blind-next-value { font-size: 20px; }
  .panel-value { font-size: 24px; }
  .break-card .panel-value { font-size: 28px; }
  .players-card .players-display { font-size: 20px; }

  .prize-list {
    flex: 0 0 auto;
    height: 100px;
    position: relative;
    min-height: 40px;
    overflow: hidden;
  }

  .mascot-area { display: none; }
  .prize-header { padding-bottom: var(--gap-sm); margin-bottom: var(--gap-sm); }
  .prize-list-inner { position: absolute; width: 100%; top: 0; left: 0; overflow: hidden; }
  .prize-item { padding: 4px 0; font-size: 16px; }
  .prize-inmoney { font-size: 16px; }

  .right-panel .info-card { padding: 10px 12px; }
  .right-panel .panel-label { font-size: 16px; margin-bottom: 4px; }
  .left-panel .panel-label { font-size: 16px; }
  .stack-card .stack-row { padding: 4px 0; }
  .stack-card .stack-label { font-size: 12px; }
  .stack-card .stack-value { font-size: 20px; }
}

/* ============================================
   レスポンシブ: 小型スマホ縦画面（max-width: 500px）
   ============================================ */
@media (max-width: 500px) and (orientation: portrait) {
  .poker-timer-app { padding: 10px; }
  .timer-grid { gap: var(--gap-sm); }
  .timer-time { font-size: 56px; }
  .blind-current { font-size: 28px; }
  .blind-current .ante-value { font-size: 28px; }
  .level-badge { font-size: 16px; padding: 6px 16px; }

  .right-panel { gap: 6px; }
  .right-panel .info-card { padding: var(--gap-sm) 10px; }
  .right-panel .panel-label { font-size: 14px; margin-bottom: 2px; }
  .left-panel .panel-label { font-size: 14px; }
  .panel-value { font-size: 20px; }
  .break-card .panel-value { font-size: 24px; }
  .players-card .players-display { font-size: 16px; }

  .stack-card .stack-row { padding: 2px 0; }
  .stack-card .stack-label { font-size: 10px; }
  .stack-card .stack-value { font-size: 16px; }
  .players-buttons .btn { width: 32px; height: 32px; font-size: 14px; }

  .left-panel { height: 210px; max-height: 210px; }
  .prize-item { font-size: 14px; }
  .prize-inmoney { font-size: 14px; }
  .prize-list { flex: 0 0 auto; height: 100px; overflow: hidden; }
  .prize-list-inner { position: absolute; width: 100%; top: 0; left: 0; will-change: transform; overflow: hidden; }

  .controls { flex-wrap: wrap; gap: 6px; }
  .btn { padding: 10px 14px; font-size: 13px; }
  .label-desktop { display: none; }
  .label-mobile { display: inline; }

  /* 設定: モバイル2行レイアウト */
  .blind-level-item {
    display: grid;
    grid-template-columns: 36px 24px 1fr 1fr 14px 32px;
    grid-template-rows: auto auto;
    row-gap: 4px;
    column-gap: 4px;
    align-items: center;
    position: relative;
    padding-right: 36px;
  }

  .blind-level-item .drag-handle { display: none; }
  .blind-level-item .move-controls {
    display: flex;
    grid-row: 1 / span 2;
    grid-column: 1;
    justify-content: flex-start;
    margin-right: -2px;
  }
  .blind-level-item .level-num { grid-row: 1; grid-column: 2; justify-self: center; min-width: 24px; }
  .blind-level-item .time-input { grid-row: 2; grid-column: 4; width: 100%; }
  .blind-level-item .sb-label, .blind-level-item .bb-label, .blind-level-item .ante-label, .blind-level-item .time-label { display: none; }
  .blind-level-item .sb-input { grid-row: 1; grid-column: 3; width: 100%; }
  .blind-level-item .bb-input { grid-row: 1; grid-column: 4; width: 100%; }
  .blind-level-item .ante-input { grid-row: 2; grid-column: 3; width: 100%; }
  .blind-level-item input { min-width: 40px; margin-right: 2px; }
  .blind-level-item .delete-level {
    position: absolute;
    right: 6px;
    top: 50%;
    transform: translateY(-50%);
    padding: 4px 6px;
    font-size: 12px;
  }
  .blind-level-item span:not(.level-num):not(.break-label):not(.break-unit) { display: none; }

  .blind-level-item.break-item {
    display: grid;
    grid-template-columns: 36px 24px 1fr 1fr 14px 32px;
    grid-template-rows: auto auto;
    gap: 6px;
    align-items: center;
    position: relative;
    padding-right: 36px;
  }
  .blind-level-item.break-item .drag-handle { display: none; }
  .blind-level-item.break-item .move-controls { display: flex; grid-column: 1; }
  .blind-level-item.break-item .break-label { grid-column: 3; grid-row: 1 / span 2; align-self: center; }
  .blind-level-item.break-item .break-time-input { grid-column: 4; grid-row: 1 / span 2; width: 100% !important; max-width: 64px; margin-right: 4px; }
  .blind-level-item.break-item .break-unit { grid-column: 5; grid-row: 1 / span 2; margin-left: 2px; align-self: center; line-height: 1; justify-self: start; position: relative; top: 2px; }
  .blind-level-item.break-item .delete-level { position: absolute; right: 6px; top: 50%; transform: translateY(-50%); }

  .prize-edit-item {
    display: grid;
    grid-template-columns: 60px 1fr 24px 1fr 24px 32px;
    grid-template-rows: auto auto;
    gap: 6px;
    align-items: center;
  }
  .prize-edit-item .drag-handle { display: none; }
  .prize-edit-item .move-controls { display: flex; grid-row: 1 / span 2; grid-column: 1; justify-content: flex-start; }
  .prize-edit-item input.rank-input { grid-row: 1; width: 100%; }
  .prize-edit-item input.rank-start { grid-column: 2; }
  .prize-edit-item input.rank-end { grid-column: 4; }
  .prize-edit-item span { font-size: 12px; white-space: nowrap; }
  .prize-edit-item span:nth-of-type(2) { grid-row: 1; grid-column: 3; justify-self: center; }
  .prize-edit-item span:nth-of-type(3) { grid-row: 1; grid-column: 5; justify-self: center; }
  .prize-edit-item input.amount-input { grid-row: 2; grid-column: 2 / span 3; width: 100%; }
  .prize-edit-item span:nth-of-type(4) { grid-row: 2; grid-column: 5; justify-self: start; }
  .prize-edit-item .delete-prize { grid-row: 1 / span 2; grid-column: 6; }
}
</style>

<!-- メインアプリ -->
<div class="poker-timer-app" id="pokerTimer" style="display: block;">
  <!-- フルスクリーン解除ボタン（モバイルフルスクリーン時のみ表示） -->
  <button class="exit-fullscreen-btn" id="btnExitFullscreen" title="フルスクリーン解除">✕</button>
  <div class="timer-grid">
    <!-- 左カラム: PRIZE -->
    <div class="left-panel">
      <div class="prize-header">
        <div class="panel-label">PRIZE</div>
        <div class="prize-inmoney" id="prizeInmoney">インマネ: <span class="prize-inmoney-value">2</span>名</div>
      </div>
      <div class="prize-list" id="prizeListContainer">
        <div class="prize-list-inner" id="prizeList">
        </div>
      </div>
      <!-- 画像領域（横画面フルスクリーン時のみ表示） -->
      <div class="mascot-area" id="mascotArea">
        <img src="data:image/jpeg;base64,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" alt="ポーカーやろうよ" class="mascot-img">
      </div>
    </div>

    <!-- 中央カラム: タイマー -->
    <div class="center-panel">
      <!-- モバイル用: フルスクリーンボタン（右上絶対配置） -->
      <button class="fullscreen-btn-mobile" id="btnFullscreenMobile" title="フルスクリーン"><i class="fas fa-expand"></i></button>

      <!-- モバイル用: Levelバッジ行 -->
      <div class="level-row">
        <div class="level-badge" id="levelBadgeMobile">Level <span id="currentLevelMobile">1</span> / <span id="totalLevelsMobile">10</span></div>
      </div>

      <!-- PC用: 単独のLevelバッジ -->
      <div class="level-badge" id="levelBadge">Level <span id="currentLevel">1</span> / <span id="totalLevels">10</span></div>

      <div class="timer-time" id="timerDisplay">10:00</div>

      <div class="progress-bar">
        <div class="progress-fill" id="progressBar" style="width: 100%"></div>
      </div>

      <div class="blind-info">
        <div class="blind-current">
          <span id="currentBlind">25 / 50</span>
          <span class="ante-value" id="anteInfo">Ante: 0</span>
        </div>
        <div class="blind-next">
          <div class="blind-next-label">Next</div>
          <span class="blind-next-value" id="nextBlind">50 / 100 <span class="next-ante">(Ante: 0)</span></span>
        </div>
      </div>

      <div class="controls">
        <button class="btn btn-primary" id="startBtn">開始</button>
        <button class="btn btn-secondary" id="btnPrev">戻る</button>
        <button class="btn btn-secondary" id="btnSkip">進む</button>
        <button class="btn btn-warning" id="btnSettings">設定</button>
      </div>
    </div>

    <!-- 右カラム -->
    <div class="right-panel">
      <!-- フルスクリーンボタン（PC用・right-panel右上） -->
      <button class="fullscreen-btn-top" id="btnFullscreen" title="フルスクリーン"><i class="fas fa-expand"></i></button>

      <!-- スペーサー -->
      <div class="right-panel-spacer"></div>

      <!-- 下部固定: NEXT BREAK + STACK + PLAYERS -->
      <div class="right-panel-bottom">
        <!-- NEXT BREAK -->
        <div class="info-card break-card" id="breakCard">
          <div class="panel-label">NEXT BREAK</div>
          <div class="panel-value large" id="nextBreakTime">--:--</div>
        </div>

        <!-- STACK -->
        <div class="info-card stack-card">
          <div class="panel-label">STACK</div>
          <div class="stack-row">
            <span class="stack-label">Avg</span>
            <span class="stack-value" id="avgStack">500</span>
          </div>
          <div class="stack-row">
            <span class="stack-label">Total</span>
            <span class="stack-value" id="totalStack">4,000</span>
          </div>
        </div>

        <!-- PLAYERS -->
        <div class="info-card players-card">
          <div class="panel-label">PLAYERS</div>
          <div class="players-display">
            <span id="remainingDisplay">8</span> / <span id="entriesDisplay">8</span>
          </div>
          <div class="players-label">Remaining / Entries</div>
          <div class="counter-controls">
            <button class="counter-btn" id="btnEntryMinus" title="エントリー減"><i class="fas fa-chevron-left"></i></button>
            <button class="counter-btn" id="btnRemainPlus" title="残り+"><i class="fas fa-chevron-up"></i></button>
            <button class="counter-btn" id="btnRemainMinus" title="残り-"><i class="fas fa-chevron-down"></i></button>
            <button class="counter-btn" id="btnEntryPlus" title="エントリー増"><i class="fas fa-chevron-right"></i></button>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- 使い方説明 -->
<div class="usage-guide">
  <h3>使い方</h3>

  <h4>基本操作</h4>
  <ul>
    <li><strong>開始/一時停止</strong>: 「開始」ボタンまたは<kbd>スペース</kbd>キーでタイマーをスタート/停止</li>
    <li><strong>戻る/進む</strong>: ブラインドレベルを前後に移動</li>
    <li><strong>設定</strong>: ブラインド構成、プライズ配分などをカスタマイズ</li>
    <li><strong>フルスクリーン</strong>: 右上のボタンで全画面表示（モバイルも対応）</li>
  </ul>

  <h4>プレイヤー数の調整</h4>
  <ul>
    <li><kbd>◀</kbd> <kbd>▶</kbd>: エントリー数を増減</li>
    <li><kbd>▲</kbd> <kbd>▼</kbd>: 残りプレイヤー数を増減</li>
    <li>キーボード: 矢印キー（←→: エントリー、↑↓: 残り人数）</li>
  </ul>

  <h4>設定のカスタマイズ</h4>
  <ul>
    <li><strong>基本設定</strong>: レベル時間、エントリー数、初期スタックを設定</li>
    <li><strong>ブラインド</strong>: SB/BB/Anteを自由に設定。ブレイクも追加可能</li>
    <li><strong>プライズ</strong>: 参加費や還元率を設定し、プライズを自動計算</li>
  </ul>

  <h4>便利な機能</h4>
  <ul>
    <li>設定は自動保存されます（ブラウザのLocalStorage）</li>
    <li>ブラインドセットを複数保存して切り替え可能</li>
    <li>ドラッグ&ドロップでブラインドレベルやプライズの並べ替えができます</li>
  </ul>
</div>

<!-- 設定モーダル -->
<div class="settings-modal" id="settingsModal" style="position:fixed!important;top:0!important;left:0!important;width:100%!important;height:100%!important;z-index:2147483647!important;">
  <div class="settings-content">
    <h2 class="settings-title">タイマー設定</h2>

    <!-- タブ -->
    <div class="settings-tabs">
      <div class="settings-tab active" id="tabBasic" data-tab="basic">基本設定</div>
      <div class="settings-tab" id="tabBlinds" data-tab="blinds">ブラインド</div>
      <div class="settings-tab" id="tabPrize" data-tab="prize">プライズ</div>
    </div>

    <!-- 基本設定パネル -->
    <div class="settings-panel active" id="panel-basic">
      <div class="setting-row">
        <div class="setting-group">
          <label class="setting-label">デフォルト時間(分)</label>
          <input type="number" class="setting-input" id="defaultLevelTime" value="20" min="1" max="60">
          <small style="color: #64748b; font-size: 11px;">レベル追加時デフォルト</small>
        </div>
        <div class="setting-group">
          <label class="setting-label">初期エントリー数</label>
          <input type="number" class="setting-input" id="initialPlayers" value="9" min="2" max="100">
        </div>
      </div>
      <div class="setting-group">
        <label class="setting-label">初期スタック</label>
        <input type="number" class="setting-input" id="initialStackInput" value="500" min="100" step="100">
      </div>
    </div>

    <!-- ブラインドパネル -->
    <div class="settings-panel" id="panel-blinds">
      <div class="setting-group">
        <label class="setting-label">ブラインドセット</label>
        <div class="blind-set-selector">
          <select id="blindSetSelect">
            <option value="default">デフォルト</option>
          </select>
          <button class="btn-save-set" id="btnSaveSet">保存</button>
          <button class="btn-delete-set" id="btnDeleteSet">削除</button>
        </div>
        <input type="text" class="setting-input" id="blindSetName" placeholder="新しいセット名" style="margin-top: 8px;">
      </div>

      <div class="setting-group">
        <label class="setting-label">
          ブラインドレベル <span class="label-desktop">(ドラッグで並べ替え可能)</span><span class="label-mobile">(ボタンで並べ替え可能)</span>
        </label>
        <div class="blind-levels" id="blindLevels"></div>
        <div class="level-actions">
          <button class="btn-add-level" id="btnAddLevel">+ レベル追加</button>
          <button class="btn-add-break" id="btnAddBreak">+ ブレイク追加</button>
        </div>
      </div>
    </div>

    <!-- プライズパネル -->
    <div class="settings-panel" id="panel-prize">
      <div class="setting-row">
        <div class="setting-group">
          <label class="setting-label">エントリー費</label>
          <input type="number" class="setting-input" id="entryFee" value="1000" min="100" step="100">
        </div>
        <div class="setting-group">
          <label class="setting-label">還元率 (%)</label>
          <input type="number" class="setting-input" id="returnRate" value="60" min="1" max="100">
        </div>
      </div>
      <div class="setting-row">
        <div class="setting-group">
          <label class="setting-label">インマネ率 (%)</label>
          <input type="number" class="setting-input" id="inMoneyRate" value="20" min="1" max="100">
        </div>
        <div class="setting-group">
          <label class="setting-label">丸め単位 (pt)</label>
          <input type="number" class="setting-input" id="roundUnit" value="100" min="1">
        </div>
      </div>
      <div class="setting-group">
        <label class="setting-label">ALPHA（上位傾斜）</label>
        <input type="number" class="setting-input" id="prizeAlpha" value="1.0" min="0.1" max="3" step="0.1">
        <small style="color: #64748b; font-size: 12px;">1.0 = 標準、大きいと上位に厚く</small>
      </div>

      <div class="setting-group">
        <label class="setting-label">
          プライズ配分 <span class="label-desktop">(ドラッグで並べ替え可能)</span><span class="label-mobile">(ボタンで並べ替え可能)</span>
        </label>
        <div class="prize-calc-row">
          <button class="btn-calc" id="btnCalcPrize">計算して反映</button>
          <small style="color: #64748b; font-size: 11px;">計算結果は下で編集可能</small>
        </div>
        <div class="prize-edit-list" id="prizeEditList">
          <!-- 計算後に表示 -->
        </div>
        <div class="prize-edit-actions">
          <button class="btn-add-prize" id="btnAddPrize">+ プライズ追加</button>
        </div>
      </div>
    </div>

    <!-- ボタン -->
    <div class="modal-buttons">
      <button class="btn btn-primary" id="btnSaveSettings">保存</button>
      <button class="btn btn-secondary" id="btnCloseSettings">閉じる</button>
    </div>
  </div>
</div>

<script>
(function() {
  // デフォルトブラインド構成（レベルごとの時間付き）
  var defaultBlinds = [
    { sb: 2, bb: 5, ante: 5, isBreak: false, levelTime: 20 },
    { sb: 5, bb: 10, ante: 10, isBreak: false, levelTime: 20 },
    { sb: 10, bb: 20, ante: 20, isBreak: false, levelTime: 20 },
    { sb: 20, bb: 40, ante: 40, isBreak: false, levelTime: 20 },
    { isBreak: true, breakTime: 10 },
    { sb: 30, bb: 60, ante: 60, isBreak: false, levelTime: 20 },
    { sb: 40, bb: 80, ante: 80, isBreak: false, levelTime: 20 },
    { sb: 50, bb: 100, ante: 100, isBreak: false, levelTime: 20 },
    { sb: 75, bb: 150, ante: 150, isBreak: false, levelTime: 20 },
    { isBreak: true, breakTime: 10 },
    { sb: 100, bb: 200, ante: 200, isBreak: false, levelTime: 20 },
    { sb: 200, bb: 400, ante: 400, isBreak: false, levelTime: 20 },
    { sb: 300, bb: 600, ante: 600, isBreak: false, levelTime: 20 },
    { sb: 400, bb: 800, ante: 800, isBreak: false, levelTime: 20 },
    { sb: 500, bb: 1000, ante: 1000, isBreak: false, levelTime: 20 }
  ];

  // プライズ設定のデフォルト
  var defaultPrizeSettings = {
    entryFee: 1000,
    returnRate: 0.60,
    inMoneyRate: 0.20,
    alpha: 1.0,
    roundUnit: 100
  };

  // 状態管理
  var state = {
    isRunning: false,
    currentLevel: 0,
    timeRemaining: 600,
    levelTime: 600,
    defaultLevelTime: 20,
    totalPlayers: 9,
    remainingPlayers: 9,
    initialStack: 500,
    blinds: JSON.parse(JSON.stringify(defaultBlinds)),
    blindSets: { 'default': JSON.parse(JSON.stringify(defaultBlinds)) },
    currentSetName: 'default',
    prizeSettings: JSON.parse(JSON.stringify(defaultPrizeSettings)),
    prizeDistribution: [],
    intervalId: null,
    scrollIntervalId: null,
    draggedItem: null,
    draggedType: null
  };

  // 要素取得ヘルパー
  function $(id) {
    return document.getElementById(id);
  }

  // レベル番号計算（ブレイクを除いた番号）
  function getLevelNumber(upToIndex) {
    var levelNum = 0;
    var limit = Math.min(upToIndex, state.blinds.length - 1);
    for (var i = 0; i <= limit; i++) {
      if (!state.blinds[i].isBreak) levelNum++;
    }
    return levelNum;
  }

  // 総レベル数計算（ブレイクを除く）
  function getTotalLevels() {
    var count = 0;
    for (var i = 0; i < state.blinds.length; i++) {
      if (!state.blinds[i].isBreak) count++;
    }
    return count;
  }

  // 現在レベルの時間を取得・設定
  function setLevelTime(levelIndex) {
    var blind = state.blinds[levelIndex];
    if (blind) {
      if (blind.isBreak) {
        state.levelTime = (blind.breakTime || 5) * 60;
      } else {
        state.levelTime = (blind.levelTime || state.defaultLevelTime) * 60;
      }
      state.timeRemaining = state.levelTime;
    }
  }

  // アプリ初期化
  function initApp() {
    var app = $('pokerTimer');
    if (app) app.style.display = 'block';
    loadSettings();
  }

  var SETTINGS_KEY = 'pokerTimerSettingsV8';

  // LocalStorageから設定を読み込み
  function loadSettings() {
    var saved = localStorage.getItem(SETTINGS_KEY);
    if (saved) {
      try {
        var settings = JSON.parse(saved);
        state.defaultLevelTime = settings.defaultLevelTime || 20;
        state.totalPlayers = settings.totalPlayers || 9;
        state.remainingPlayers = settings.totalPlayers || 9;
        state.initialStack = settings.initialStack || 500;
        state.blinds = settings.blinds || JSON.parse(JSON.stringify(defaultBlinds));
        state.blindSets = settings.blindSets || { 'default': JSON.parse(JSON.stringify(defaultBlinds)) };
        state.currentSetName = settings.currentSetName || 'default';
        state.prizeSettings = settings.prizeSettings || JSON.parse(JSON.stringify(defaultPrizeSettings));
        state.prizeDistribution = settings.prizeDistribution || [];
      } catch(e) {
        console.log('Settings load error:', e);
      }
    }

    // 現在レベルの時間を設定（保存データの有無に関わらず実行）
    setLevelTime(state.currentLevel);

    updateBlindSetSelector();
    updateDisplay();
    updatePrizeDisplay();
    startPrizeAutoScroll();
  }

  // 設定を保存
  function saveSettings() {
    var defaultLevelTime = parseInt($('defaultLevelTime').value) || 20;
    var totalPlayers = parseInt($('initialPlayers').value) || 9;
    var initialStack = parseInt($('initialStackInput').value) || 500;

    // ブラインドレベルを収集
    var blindLevels = document.querySelectorAll('.blind-level-item');
    var blinds = [];
    for (var i = 0; i < blindLevels.length; i++) {
      var item = blindLevels[i];
      if (item.classList.contains('break-item')) {
        var breakTime = parseInt(item.querySelector('.break-time-input').value) || 5;
        blinds.push({ isBreak: true, breakTime: breakTime });
      } else {
        var time = parseInt(item.querySelector('.time-input').value) || defaultLevelTime;
        var sb = parseInt(item.querySelector('.sb-input').value) || 0;
        var bb = parseInt(item.querySelector('.bb-input').value) || 0;
        var ante = parseInt(item.querySelector('.ante-input').value) || 0;
        blinds.push({ sb: sb, bb: bb, ante: ante, isBreak: false, levelTime: time });
      }
    }

    // プライズ設定
    var prizeSettings = {
      entryFee: parseInt($('entryFee').value) || 1000,
      returnRate: (parseInt($('returnRate').value) || 60) / 100,
      inMoneyRate: (parseInt($('inMoneyRate').value) || 20) / 100,
      alpha: parseFloat($('prizeAlpha').value) || 1.0,
      roundUnit: parseInt($('roundUnit').value) || 100
    };

    // プライズ配分を収集（範囲指定対応）
    var prizeItems = document.querySelectorAll('.prize-edit-item');
    var prizeDistribution = [];
    for (var j = 0; j < prizeItems.length; j++) {
      var startRank = parseInt(prizeItems[j].querySelector('.rank-start').value) || 1;
      var endRank = parseInt(prizeItems[j].querySelector('.rank-end').value) || startRank;
      var amount = parseInt(prizeItems[j].querySelector('.amount-input').value) || 0;
      prizeDistribution.push({ startRank: startRank, endRank: endRank, amount: amount });
    }

    state.defaultLevelTime = defaultLevelTime;
    state.totalPlayers = totalPlayers;
    state.remainingPlayers = totalPlayers;
    state.initialStack = initialStack;
    state.blinds = blinds.length > 0 ? blinds : JSON.parse(JSON.stringify(defaultBlinds));
    state.prizeSettings = prizeSettings;
    state.prizeDistribution = prizeDistribution;
    state.currentLevel = 0;

    // 現在レベルの時間を設定
    setLevelTime(0);

    localStorage.setItem(SETTINGS_KEY, JSON.stringify({
      defaultLevelTime: defaultLevelTime,
      totalPlayers: totalPlayers,
      initialStack: initialStack,
      blinds: state.blinds,
      blindSets: state.blindSets,
      currentSetName: state.currentSetName,
      prizeSettings: state.prizeSettings,
      prizeDistribution: state.prizeDistribution
    }));

    updateDisplay();
    updatePrizeDisplay();
    closeSettings();
  }

  // ブラインドセット管理
  function saveBlindSet() {
    var name = $('blindSetName').value.trim();
    if (!name) {
      alert('セット名を入力してください');
      return;
    }

    var blindLevels = document.querySelectorAll('.blind-level-item');
    var blinds = [];
    var defaultTime = parseInt($('defaultLevelTime').value) || 20;
    for (var i = 0; i < blindLevels.length; i++) {
      var item = blindLevels[i];
      if (item.classList.contains('break-item')) {
        var breakTime = parseInt(item.querySelector('.break-time-input').value) || 5;
        blinds.push({ isBreak: true, breakTime: breakTime });
      } else {
        var time = parseInt(item.querySelector('.time-input').value) || defaultTime;
        var sb = parseInt(item.querySelector('.sb-input').value) || 0;
        var bb = parseInt(item.querySelector('.bb-input').value) || 0;
        var ante = parseInt(item.querySelector('.ante-input').value) || 0;
        blinds.push({ sb: sb, bb: bb, ante: ante, isBreak: false, levelTime: time });
      }
    }

    state.blindSets[name] = blinds;
    state.currentSetName = name;
    updateBlindSetSelector();
    $('blindSetSelect').value = name;
    $('blindSetName').value = '';
    alert('セット「' + name + '」を保存しました');
  }

  function loadBlindSet() {
    var name = $('blindSetSelect').value;
    if (state.blindSets[name]) {
      state.blinds = JSON.parse(JSON.stringify(state.blindSets[name]));
      state.currentSetName = name;
      renderBlindLevels();
    }
  }

  function deleteBlindSet() {
    var name = $('blindSetSelect').value;
    if (name === 'default') {
      alert('デフォルトセットは削除できません');
      return;
    }
    if (confirm('セット「' + name + '」を削除しますか？')) {
      delete state.blindSets[name];
      state.currentSetName = 'default';
      state.blinds = JSON.parse(JSON.stringify(state.blindSets['default']));
      updateBlindSetSelector();
      renderBlindLevels();
    }
  }

  function updateBlindSetSelector() {
    var select = $('blindSetSelect');
    if (!select) return;
    select.innerHTML = '';
    var names = Object.keys(state.blindSets);
    for (var i = 0; i < names.length; i++) {
      var name = names[i];
      var option = document.createElement('option');
      option.value = name;
      option.textContent = name === 'default' ? 'デフォルト' : name;
      if (name === state.currentSetName) option.selected = true;
      select.appendChild(option);
    }
  }

  // タブ切り替え
  function switchTab(tabName) {
    var tabs = document.querySelectorAll('.settings-tab');
    var panels = document.querySelectorAll('.settings-panel');
    for (var i = 0; i < tabs.length; i++) {
      tabs[i].classList.remove('active');
    }
    for (var j = 0; j < panels.length; j++) {
      panels[j].classList.remove('active');
    }
    var activeTab = document.querySelector('.settings-tab[data-tab="' + tabName + '"]');
    if (activeTab) activeTab.classList.add('active');
    var activePanel = $('panel-' + tabName);
    if (activePanel) activePanel.classList.add('active');
  }

  // エントリー数調整
  function adjustEntries(delta) {
    var oldTotal = state.totalPlayers;
    state.totalPlayers = Math.max(1, state.totalPlayers + delta);
    // エントリー増加時は残り人数も同じだけ増やす
    if (delta > 0) {
      state.remainingPlayers = Math.min(state.totalPlayers, state.remainingPlayers + delta);
    }
    // エントリー減少時は残り人数が超えないように調整
    if (state.remainingPlayers > state.totalPlayers) {
      state.remainingPlayers = state.totalPlayers;
    }
    updateDisplay();
    updatePrizeDisplay();
  }

  // 残りプレイヤー数調整
  function adjustRemaining(delta) {
    state.remainingPlayers = Math.max(1, Math.min(state.totalPlayers, state.remainingPlayers + delta));
    updateDisplay();
  }

  // タイマー開始/停止
  function toggleTimer() {
    state.isRunning = !state.isRunning;
    var startBtn = $('startBtn');

    if (state.isRunning) {
      initAudio(); // iOS対応: ユーザー操作時にAudioContextを初期化
      if (startBtn) startBtn.textContent = '停止';
      state.intervalId = setInterval(tick, 1000);
    } else {
      clearInterval(state.intervalId);
      if (startBtn) startBtn.textContent = '開始';
    }
  }

  function tick() {
    state.timeRemaining--;

    if (state.timeRemaining <= 0) {
      playSound();
      nextLevel();
    }

    updateDisplay();
  }

  function nextLevel() {
    if (state.currentLevel < state.blinds.length - 1) {
      state.currentLevel++;
      setLevelTime(state.currentLevel);
      // 音は tick() から呼ばれる時のみ再生（進むボタンでは鳴らさない）
    } else {
      clearInterval(state.intervalId);
      state.isRunning = false;
      var startBtn = $('startBtn');
      if (startBtn) startBtn.textContent = '開始';
    }
    updateDisplay();
  }

  function prevLevel() {
    if (state.currentLevel > 0) {
      state.currentLevel--;
      setLevelTime(state.currentLevel);
      updateDisplay();
    }
  }

  function skipLevel() {
    nextLevel();
  }

  // 時間フォーマット（MM:SS）
  function formatTime(seconds) {
    var m = Math.floor(seconds / 60);
    var s = seconds % 60;
    return String(m).padStart(2, '0') + ':' + String(s).padStart(2, '0');
  }

  // レベルバッジ更新
  function updateLevelBadge(blind) {
    var lb = $('levelBadge');
    var lbMobile = $('levelBadgeMobile');

    if (blind.isBreak) {
      if (lb) { lb.textContent = 'BREAK'; lb.classList.add('break'); }
      if (lbMobile) { lbMobile.textContent = 'BREAK'; lbMobile.classList.add('break'); }
    } else {
      var levelNum = getLevelNumber(state.currentLevel);
      var totalLevelsCount = getTotalLevels();
      var html = 'Level <span id="currentLevel">' + levelNum + '</span> / <span id="totalLevels">' + totalLevelsCount + '</span>';
      if (lb) { lb.innerHTML = html; lb.classList.remove('break'); }
      if (lbMobile) {
        lbMobile.innerHTML = html.replace('currentLevel', 'currentLevelMobile').replace('totalLevels', 'totalLevelsMobile');
        lbMobile.classList.remove('break');
      }
    }
  }

  // 表示更新
  function updateDisplay() {
    var td = $('timerDisplay');
    if (td) {
      td.textContent = formatTime(state.timeRemaining);
      td.classList.toggle('warning', state.timeRemaining <= 10);
    }

    // プログレスバー
    var progress = state.levelTime > 0 ? (state.timeRemaining / state.levelTime) * 100 : 0;
    var pb = $('progressBar');
    if (pb) pb.style.width = progress + '%';

    // ブラインド情報
    var blind = state.blinds[state.currentLevel];
    var cb = $('currentBlind');
    var ai = $('anteInfo');
    var nb = $('nextBlind');

    if (blind) {
      updateLevelBadge(blind);
      if (blind.isBreak) {
        if (cb) cb.textContent = 'BREAK TIME';
        if (ai) ai.textContent = '';
      } else {
        if (cb) cb.textContent = blind.sb + ' / ' + blind.bb;
        if (ai) ai.textContent = 'Ante: ' + blind.ante;
      }
    }

    // 次のブラインド（Ante付き）
    var nextBlindData = null;
    for (var k = state.currentLevel + 1; k < state.blinds.length; k++) {
      if (!state.blinds[k].isBreak) {
        nextBlindData = state.blinds[k];
        break;
      }
    }
    if (nb) {
      nb.innerHTML = nextBlindData
        ? nextBlindData.sb + ' / ' + nextBlindData.bb + ' <span class="next-ante">(Ante: ' + nextBlindData.ante + ')</span>'
        : '--';
    }

    // スタック表示
    var totalStack = state.initialStack * state.totalPlayers;
    var avgStack = state.remainingPlayers > 0 ? Math.round(totalStack / state.remainingPlayers) : 0;
    var ts = $('totalStack');
    var as = $('avgStack');
    if (ts) ts.textContent = totalStack.toLocaleString();
    if (as) as.textContent = avgStack.toLocaleString();

    // プレイヤー表示
    var rd = $('remainingDisplay');
    var ed = $('entriesDisplay');
    if (rd) rd.textContent = state.remainingPlayers;
    if (ed) ed.textContent = state.totalPlayers;

    // NEXT BREAK IN 更新
    updateNextBreak();
  }

  // NEXT BREAK IN 表示更新
  function updateNextBreak() {
    var breakCard = $('breakCard');
    if (!breakCard) return;

    var timeToBreak = state.timeRemaining;
    var foundBreak = false;

    for (var i = state.currentLevel + 1; i < state.blinds.length; i++) {
      if (state.blinds[i].isBreak) {
        foundBreak = true;
        break;
      } else {
        timeToBreak += (state.blinds[i].levelTime || state.defaultLevelTime) * 60;
      }
    }

    var nbt = $('nextBreakTime');
    var showBreak = foundBreak ? !state.blinds[state.currentLevel].isBreak : false;

    breakCard.classList.toggle('no-break', !showBreak);
    if (nbt) nbt.textContent = showBreak ? formatTime(timeToBreak) : '--:--';
  }

  // プライズ計算・表示（範囲指定対応）
  function updatePrizeDisplay() {
    var ps = state.prizeSettings;
    var N = state.totalPlayers;
    var P = N * ps.entryFee * ps.returnRate;

    // インマネ人数計算
    var totalInMoney = 0;
    if (state.prizeDistribution.length > 0) {
      for (var x = 0; x < state.prizeDistribution.length; x++) {
        var pd = state.prizeDistribution[x];
        totalInMoney = Math.max(totalInMoney, pd.endRank);
      }
    } else {
      totalInMoney = Math.ceil(N * ps.inMoneyRate);
      if (totalInMoney < 1) totalInMoney = 1;
    }

    var pi = $('prizeInmoney');
    if (pi) pi.innerHTML = 'インマネ: <span class="prize-inmoney-value">' + totalInMoney + '</span>名';

    var listEl = $('prizeList');
    if (listEl) {
      listEl.innerHTML = '';
      listEl.removeAttribute('data-duplicated');  // 複製フラグをリセット
      listEl.style.transform = '';  // スクロール位置をリセット

      if (state.prizeDistribution.length > 0) {
        // 手動設定を使用
        for (var k = 0; k < state.prizeDistribution.length; k++) {
          var p = state.prizeDistribution[k];
          var item = document.createElement('div');
          item.className = 'prize-item';

          var rankText = '';
          if (p.startRank === p.endRank) {
            rankText = p.startRank + '位';
          } else {
            rankText = p.startRank + '～' + p.endRank + '位';
          }

          item.innerHTML = '<span class="prize-rank">' + rankText + '</span><span class="prize-amount">' + p.amount.toLocaleString() + ' pt</span>';
          listEl.appendChild(item);
        }
      } else {
        // 自動計算
        var M = Math.ceil(N * ps.inMoneyRate);
        if (M < 1) M = 1;

        var W = 0;
        for (var i = 1; i <= M; i++) {
          W += 1 / Math.pow(i, ps.alpha);
        }

        var prizes = [];
        var total = 0;
        for (var j = 1; j <= M; j++) {
          var weight = 1 / Math.pow(j, ps.alpha);
          var raw = P * weight / W;
          var rounded = Math.floor(raw / ps.roundUnit) * ps.roundUnit;
          prizes.push(rounded);
          total += rounded;
        }

        var diff = Math.round(P) - total;
        if (prizes.length > 0) {
          prizes[0] += diff;
        }

        for (var m = 0; m < prizes.length; m++) {
          var item2 = document.createElement('div');
          item2.className = 'prize-item';
          item2.innerHTML = '<span class="prize-rank">' + (m + 1) + '位</span><span class="prize-amount">' + prizes[m].toLocaleString() + ' pt</span>';
          listEl.appendChild(item2);
        }
      }
    }

    setupPrizeScroll();
  }

  function setupPrizeScroll() {
    var container = $('prizeListContainer');
    var inner = $('prizeList');

    if (!container || !inner) return;

    var containerHeight = container.clientHeight;
    var innerHeight = inner.scrollHeight;
    if (!containerHeight || !innerHeight) return;

    if (innerHeight > containerHeight + 1) {
      inner.style.animation = 'none';
      startPrizeAutoScroll();
    } else {
      if (state.scrollIntervalId) {
        clearInterval(state.scrollIntervalId);
        state.scrollIntervalId = null;
      }
      if (inner.getAttribute('data-duplicated')) {
        if (inner.getAttribute('data-original-html')) {
          inner.innerHTML = inner.getAttribute('data-original-html');
          inner.removeAttribute('data-duplicated');
          inner.removeAttribute('data-original-html');
        }
      }
      inner.style.top = '0';
      inner.style.transform = '';
    }
  }

  // 無限ループスクロール（おしりと頭の間に5行分の隙間）
  function startPrizeAutoScroll() {
    if (state.scrollIntervalId) {
      clearInterval(state.scrollIntervalId);
    }

    var container = $('prizeListContainer');
    var inner = $('prizeList');

    if (!inner || !container) return;

    // 既に複製されている場合は元に戻してから高さを計算
    if (inner.getAttribute('data-duplicated')) {
      var originalHtml = inner.getAttribute('data-original-html');
      if (originalHtml) {
        inner.innerHTML = originalHtml;
        inner.removeAttribute('data-duplicated');
        inner.removeAttribute('data-original-html');
        inner.removeAttribute('data-loop-height');
        inner.style.transform = '';
      }
    }

    var containerHeight = container.clientHeight;
    var innerHeight = inner.scrollHeight;

    // コンテンツがコンテナに収まる場合はスクロールしない
    if (!containerHeight || !innerHeight || innerHeight <= containerHeight) return;

    // コンテンツを複製して継ぎ目なくループ（おしりと頭の間に隙間）
    var originalContent = inner.innerHTML;
    var originalHeight = inner.scrollHeight;
    var spacerHeight = 80;
    if (!inner.getAttribute('data-duplicated')) {
      var spacer = '<div style="height: ' + spacerHeight + 'px;"></div>';
      inner.setAttribute('data-original-html', originalContent);
      inner.setAttribute('data-loop-height', String(originalHeight + spacerHeight));
      inner.innerHTML = originalContent + spacer + originalContent;
      inner.setAttribute('data-duplicated', 'true');
    }

    var scrollPos = 0;
    var loopHeight = parseInt(inner.getAttribute('data-loop-height'), 10);
    if (!loopHeight) loopHeight = inner.scrollHeight / 2;

    state.scrollIntervalId = setInterval(function() {
      scrollPos += 1;

      if (scrollPos >= loopHeight) {
        scrollPos -= loopHeight;  // 位置を維持したままループ
      }

      inner.style.transform = 'translate3d(0,-' + scrollPos + 'px,0)';
    }, 50);
  }

  // プライズ計算（設定画面用）
  function calculatePrizes() {
    var ps = {
      entryFee: parseInt($('entryFee').value) || 1000,
      returnRate: (parseInt($('returnRate').value) || 60) / 100,
      inMoneyRate: (parseInt($('inMoneyRate').value) || 20) / 100,
      alpha: parseFloat($('prizeAlpha').value) || 1.0,
      roundUnit: parseInt($('roundUnit').value) || 100
    };

    var N = state.totalPlayers;
    var P = N * ps.entryFee * ps.returnRate;
    var M = Math.ceil(N * ps.inMoneyRate);
    if (M < 1) M = 1;

    var W = 0;
    for (var i = 1; i <= M; i++) {
      W += 1 / Math.pow(i, ps.alpha);
    }

    var prizes = [];
    var total = 0;
    for (var j = 1; j <= M; j++) {
      var weight = 1 / Math.pow(j, ps.alpha);
      var raw = P * weight / W;
      var rounded = Math.floor(raw / ps.roundUnit) * ps.roundUnit;
      prizes.push({ startRank: j, endRank: j, amount: rounded });
      total += rounded;
    }

    var diff = Math.round(P) - total;
    if (prizes.length > 0) {
      prizes[0].amount += diff;
    }

    renderPrizeEditList(prizes);
  }

  function renderPrizeEditList(prizes) {
    var container = $('prizeEditList');
    if (!container) return;
    container.innerHTML = '';

    for (var i = 0; i < prizes.length; i++) {
      var p = prizes[i];
      var item = document.createElement('div');
      item.className = 'prize-edit-item';
      item.setAttribute('data-index', i);
      item.setAttribute('draggable', 'true');
      item.innerHTML =
        '<div class="move-controls">' +
          '<button class="move-btn move-up" title="上へ">&uarr;</button>' +
          '<button class="move-btn move-down" title="下へ">&darr;</button>' +
        '</div>' +
        '<span class="drag-handle">&equiv;</span>' +
        '<input type="number" class="rank-input rank-start" value="' + p.startRank + '" min="1">' +
        '<span>位 ～</span>' +
        '<input type="number" class="rank-input rank-end" value="' + p.endRank + '" min="1">' +
        '<span>位:</span>' +
        '<input type="number" class="amount-input" value="' + p.amount + '" min="0" step="100">' +
        '<span>pt</span>' +
        '<button class="delete-prize">&times;</button>';
      container.appendChild(item);
    }

    setupPrizeDragAndDrop();
    setupPrizeMoveButtons();
    setupPrizeDeleteButtons();
  }

  function setupPrizeMoveButtons() {
    var container = $('prizeEditList');
    if (!container) return;
    var moveUpButtons = container.querySelectorAll('.move-up');
    for (var i = 0; i < moveUpButtons.length; i++) {
      moveUpButtons[i].addEventListener('click', function(e) {
        var idx = parseInt(e.target.closest('.prize-edit-item').getAttribute('data-index'));
        syncPrizeDistributionFromInputs();
        if (idx > 0) {
          var temp = state.prizeDistribution[idx - 1];
          state.prizeDistribution[idx - 1] = state.prizeDistribution[idx];
          state.prizeDistribution[idx] = temp;
          renderPrizeEditList(state.prizeDistribution);
        }
      });
    }
    var moveDownButtons = container.querySelectorAll('.move-down');
    for (var j = 0; j < moveDownButtons.length; j++) {
      moveDownButtons[j].addEventListener('click', function(e) {
        var idx = parseInt(e.target.closest('.prize-edit-item').getAttribute('data-index'));
        syncPrizeDistributionFromInputs();
        if (idx < state.prizeDistribution.length - 1) {
          var temp = state.prizeDistribution[idx + 1];
          state.prizeDistribution[idx + 1] = state.prizeDistribution[idx];
          state.prizeDistribution[idx] = temp;
          renderPrizeEditList(state.prizeDistribution);
        }
      });
    }
  }

  function syncPrizeDistributionFromInputs() {
    var items = document.querySelectorAll('.prize-edit-item');
    var prizeDistribution = [];
    for (var i = 0; i < items.length; i++) {
      var item = items[i];
      var startRank = parseInt(item.querySelector('.rank-start').value) || 1;
      var endRank = parseInt(item.querySelector('.rank-end').value) || startRank;
      var amount = parseInt(item.querySelector('.amount-input').value) || 0;
      prizeDistribution.push({ startRank: startRank, endRank: endRank, amount: amount });
    }
    state.prizeDistribution = prizeDistribution;
  }

  function setupPrizeDeleteButtons() {
    var container = $('prizeEditList');
    if (!container) return;
    var deleteButtons = container.querySelectorAll('.delete-prize');
    for (var j = 0; j < deleteButtons.length; j++) {
      deleteButtons[j].addEventListener('click', function(e) {
        var idx = parseInt(e.target.closest('.prize-edit-item').getAttribute('data-index'));
        removePrizeRow(idx);
      });
    }
  }

  // プライズのドラッグ&ドロップ
  function setupPrizeDragAndDrop() {
    var container = $('prizeEditList');
    if (!container) return;

    var items = container.querySelectorAll('.prize-edit-item');
    for (var i = 0; i < items.length; i++) {
      items[i].addEventListener('dragstart', handlePrizeDragStart);
      items[i].addEventListener('dragend', handlePrizeDragEnd);
      items[i].addEventListener('dragover', handlePrizeDragOver);
      items[i].addEventListener('dragleave', handlePrizeDragLeave);
      items[i].addEventListener('drop', handlePrizeDrop);
    }
  }

  function handlePrizeDragStart(e) {
    state.draggedItem = this;
    state.draggedType = 'prize';
    this.classList.add('dragging');
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/plain', this.getAttribute('data-index'));
  }

  function handlePrizeDragEnd(e) {
    this.classList.remove('dragging');
    var items = document.querySelectorAll('.prize-edit-item');
    for (var i = 0; i < items.length; i++) {
      items[i].classList.remove('drag-over');
    }
    state.draggedItem = null;
    state.draggedType = null;
  }

  function handlePrizeDragOver(e) {
    if (e.preventDefault) e.preventDefault();
    if (state.draggedType !== 'prize') return;
    e.dataTransfer.dropEffect = 'move';
    this.classList.add('drag-over');
    return false;
  }

  function handlePrizeDragLeave(e) {
    this.classList.remove('drag-over');
  }

  function handlePrizeDrop(e) {
    if (e.stopPropagation) e.stopPropagation();
    if (state.draggedType !== 'prize') return;

    this.classList.remove('drag-over');

    if (state.draggedItem !== this) {
      var container = $('prizeEditList');
      var items = Array.from(container.querySelectorAll('.prize-edit-item'));
      var fromIndex = items.indexOf(state.draggedItem);
      var toIndex = items.indexOf(this);

      if (fromIndex < toIndex) {
        container.insertBefore(state.draggedItem, this.nextSibling);
      } else {
        container.insertBefore(state.draggedItem, this);
      }

      var newItems = container.querySelectorAll('.prize-edit-item');
      for (var i = 0; i < newItems.length; i++) {
        newItems[i].setAttribute('data-index', i);
      }
    }
    return false;
  }

  function addPrizeRow() {
    var container = $('prizeEditList');
    if (!container) return;

    var items = container.querySelectorAll('.prize-edit-item');
    var nextRank = 1;
    if (items.length > 0) {
      var lastItem = items[items.length - 1];
      var lastEnd = parseInt(lastItem.querySelector('.rank-end').value) || 0;
      nextRank = lastEnd + 1;
    }

    state.prizeDistribution.push({ startRank: nextRank, endRank: nextRank, amount: 0 });

    var newIndex = items.length;
    var item = document.createElement('div');
    item.className = 'prize-edit-item';
    item.setAttribute('data-index', newIndex);
    item.setAttribute('draggable', 'true');
    item.innerHTML =
      '<div class="move-controls">' +
        '<button class="move-btn move-up" title="上へ">&uarr;</button>' +
        '<button class="move-btn move-down" title="下へ">&darr;</button>' +
      '</div>' +
      '<span class="drag-handle">&equiv;</span>' +
      '<input type="number" class="rank-input rank-start" value="' + nextRank + '" min="1">' +
      '<span>位 ～</span>' +
      '<input type="number" class="rank-input rank-end" value="' + nextRank + '" min="1">' +
      '<span>位:</span>' +
      '<input type="number" class="amount-input" value="0" min="0" step="100">' +
      '<span>pt</span>' +
      '<button class="delete-prize">&times;</button>';

    container.appendChild(item);

    // イベント設定
    item.addEventListener('dragstart', handlePrizeDragStart);
    item.addEventListener('dragend', handlePrizeDragEnd);
    item.addEventListener('dragover', handlePrizeDragOver);
    item.addEventListener('dragleave', handlePrizeDragLeave);
    item.addEventListener('drop', handlePrizeDrop);

    item.querySelector('.delete-prize').addEventListener('click', function() {
      removePrizeRow(newIndex);
    });
  }

  function removePrizeRow(index) {
    var container = $('prizeEditList');
    if (!container) return;

    var items = container.querySelectorAll('.prize-edit-item');
    if (items.length <= 1) return;
    if (!items[index]) return;

    items[index].remove();
    // インデックスを再設定
    var newItems = container.querySelectorAll('.prize-edit-item');
    for (var i = 0; i < newItems.length; i++) {
      newItems[i].setAttribute('data-index', i);
    }
  }

  // 設定モーダル
  function openSettings() {
    var modal = $('settingsModal');
    if (modal) modal.classList.add('active');
    document.body.classList.add('modal-open');

    var dlt = $('defaultLevelTime');
    if (dlt) dlt.value = state.defaultLevelTime;
    var ip = $('initialPlayers');
    if (ip) ip.value = state.totalPlayers;
    var isi = $('initialStackInput');
    if (isi) isi.value = state.initialStack;

    var ef = $('entryFee');
    if (ef) ef.value = state.prizeSettings.entryFee;
    var rr = $('returnRate');
    if (rr) rr.value = state.prizeSettings.returnRate * 100;
    var imr = $('inMoneyRate');
    if (imr) imr.value = state.prizeSettings.inMoneyRate * 100;
    var pa = $('prizeAlpha');
    if (pa) pa.value = state.prizeSettings.alpha;
    var ru = $('roundUnit');
    if (ru) ru.value = state.prizeSettings.roundUnit;

    updateBlindSetSelector();
    renderBlindLevels();

    // プライズ編集リストを表示
    if (state.prizeDistribution.length > 0) {
      renderPrizeEditList(state.prizeDistribution);
    } else {
      $('prizeEditList').innerHTML = '<div style="padding: 10px; color: #64748b; text-align: center;">「計算して反映」を押すか、「+ プライズ追加」で手動追加</div>';
    }
  }

  function closeSettings() {
    var modal = $('settingsModal');
    if (modal) modal.classList.remove('active');
    document.body.classList.remove('modal-open');
  }

  function renderBlindLevels() {
    var container = $('blindLevels');
    if (!container) return;
    container.innerHTML = '';

    var levelNum = 0;
    for (var index = 0; index < state.blinds.length; index++) {
      var blind = state.blinds[index];
      var item = document.createElement('div');

      if (blind.isBreak) {
        item.className = 'blind-level-item break-item';
        item.setAttribute('data-index', index);
        item.setAttribute('draggable', 'true');
        item.innerHTML =
          '<div class="move-controls">' +
            '<button class="move-btn move-up" title="上へ">&uarr;</button>' +
            '<button class="move-btn move-down" title="下へ">&darr;</button>' +
          '</div>' +
          '<span class="drag-handle">&equiv;</span>' +
          '<span class="break-label">BREAK</span>' +
          '<input type="number" class="break-time-input" value="' + (blind.breakTime || 5) + '" min="1" style="width: 50px;"> <span class="break-unit">分</span>' +
          '<button class="delete-level">&times;</button>';
      } else {
        levelNum++;
        item.className = 'blind-level-item';
        item.setAttribute('data-index', index);
        item.setAttribute('draggable', 'true');
        item.innerHTML =
          '<div class="move-controls">' +
            '<button class="move-btn move-up" title="上へ">&uarr;</button>' +
            '<button class="move-btn move-down" title="下へ">&darr;</button>' +
          '</div>' +
          '<span class="drag-handle">&equiv;</span>' +
          '<span class="level-num">' + levelNum + '</span>' +
          '<input type="number" class="time-input" value="' + (blind.levelTime || state.defaultLevelTime) + '" min="1" title="時間(分)" placeholder="時間">' +
          '<span class="time-label">分</span>' +
          '<input type="number" class="sb-input" value="' + blind.sb + '" placeholder="SB">' +
          '<span class="sb-label">SB</span>' +
          '<input type="number" class="bb-input" value="' + blind.bb + '" placeholder="BB">' +
          '<span class="bb-label">BB</span>' +
          '<input type="number" class="ante-input" value="' + blind.ante + '" placeholder="Ante">' +
          '<span class="ante-label">Ante</span>' +
          '<button class="delete-level">&times;</button>';
      }
      container.appendChild(item);
    }

    setupBlindDragAndDrop();
    setupBlindMoveButtons();
    setupBlindDeleteButtons();
  }

  function setupBlindMoveButtons() {
    var container = $('blindLevels');
    if (!container) return;
    var moveUpButtons = container.querySelectorAll('.move-up');
    for (var i = 0; i < moveUpButtons.length; i++) {
      moveUpButtons[i].addEventListener('click', function(e) {
        var idx = parseInt(e.target.closest('.blind-level-item').getAttribute('data-index'));
        if (idx > 0) {
          var temp = state.blinds[idx - 1];
          state.blinds[idx - 1] = state.blinds[idx];
          state.blinds[idx] = temp;
          renderBlindLevels();
        }
      });
    }
    var moveDownButtons = container.querySelectorAll('.move-down');
    for (var j = 0; j < moveDownButtons.length; j++) {
      moveDownButtons[j].addEventListener('click', function(e) {
        var idx = parseInt(e.target.closest('.blind-level-item').getAttribute('data-index'));
        if (idx < state.blinds.length - 1) {
          var temp = state.blinds[idx + 1];
          state.blinds[idx + 1] = state.blinds[idx];
          state.blinds[idx] = temp;
          renderBlindLevels();
        }
      });
    }
  }

  function setupBlindDeleteButtons() {
    var container = $('blindLevels');
    if (!container) return;
    var deleteButtons = container.querySelectorAll('.delete-level');
    for (var i = 0; i < deleteButtons.length; i++) {
      deleteButtons[i].addEventListener('click', function(e) {
        var idx = parseInt(e.target.closest('.blind-level-item').getAttribute('data-index'));
        deleteBlindLevel(idx);
      });
    }
  }

  // ブラインドのドラッグ&ドロップ
  function setupBlindDragAndDrop() {
    var container = $('blindLevels');
    if (!container) return;

    var items = container.querySelectorAll('.blind-level-item');
    for (var i = 0; i < items.length; i++) {
      items[i].addEventListener('dragstart', handleBlindDragStart);
      items[i].addEventListener('dragend', handleBlindDragEnd);
      items[i].addEventListener('dragover', handleBlindDragOver);
      items[i].addEventListener('dragleave', handleBlindDragLeave);
      items[i].addEventListener('drop', handleBlindDrop);
    }
  }

  function handleBlindDragStart(e) {
    state.draggedItem = this;
    state.draggedType = 'blind';
    this.classList.add('dragging');
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/plain', this.getAttribute('data-index'));
  }

  function handleBlindDragEnd(e) {
    this.classList.remove('dragging');
    var items = document.querySelectorAll('.blind-level-item');
    for (var i = 0; i < items.length; i++) {
      items[i].classList.remove('drag-over');
    }
    state.draggedItem = null;
    state.draggedType = null;
  }

  function handleBlindDragOver(e) {
    if (e.preventDefault) e.preventDefault();
    if (state.draggedType !== 'blind') return;
    e.dataTransfer.dropEffect = 'move';
    this.classList.add('drag-over');
    return false;
  }

  function handleBlindDragLeave(e) {
    this.classList.remove('drag-over');
  }

  function handleBlindDrop(e) {
    if (e.stopPropagation) e.stopPropagation();
    if (state.draggedType !== 'blind') return;

    this.classList.remove('drag-over');

    if (state.draggedItem !== this) {
      var container = $('blindLevels');
      var items = Array.from(container.querySelectorAll('.blind-level-item'));
      var fromIndex = items.indexOf(state.draggedItem);
      var toIndex = items.indexOf(this);

      var movedBlind = state.blinds.splice(fromIndex, 1)[0];
      state.blinds.splice(toIndex, 0, movedBlind);
      renderBlindLevels();
    }
    return false;
  }

  function addBlindLevel() {
    var lastBlind = null;
    for (var i = state.blinds.length - 1; i >= 0; i--) {
      if (!state.blinds[i].isBreak) {
        lastBlind = state.blinds[i];
        break;
      }
    }
    if (!lastBlind) lastBlind = { sb: 0, bb: 0, ante: 0, levelTime: state.defaultLevelTime };
    state.blinds.push({
      sb: lastBlind.bb,
      bb: lastBlind.bb * 2,
      ante: lastBlind.ante,
      isBreak: false,
      levelTime: state.defaultLevelTime
    });
    renderBlindLevels();
  }

  function addBreakLevel() {
    state.blinds.push({ isBreak: true, breakTime: 5 });
    renderBlindLevels();
  }

  function deleteBlindLevel(index) {
    if (state.blinds.length > 1) {
      state.blinds.splice(index, 1);
      renderBlindLevels();
    }
  }

  // 無音のダミー音声データ（iOS/iPad用アンロック専用）
  // 44.1kHz, 16bit, モノラル, 0.01秒の完全無音WAV
  var SILENT_DATA = 'UklGRiQAAABXQVZFZm10IBAAAAABAAEARKwAAIhYAQACABAAZGF0YQAAAAA=';

  // チャイム音（ド→ミ→ソ→ド）のBase64データ
  var CHIME_DATA = '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AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==';

  // audio要素（iOS/iPad対応）
  var chimeAudio = null;
  var audioUnlocked = false;

  function initAudio() {
    if (!chimeAudio) {
      chimeAudio = document.createElement('audio');
      chimeAudio.setAttribute('playsinline', '');
      chimeAudio.setAttribute('webkit-playsinline', '');
      chimeAudio.src = 'data:audio/wav;base64,' + CHIME_DATA;
    }

    // iOS/iPad Safari対応: 無音データを再生してオーディオをアンロック
    if (!audioUnlocked) {
      var silentAudio = new Audio('data:audio/wav;base64,' + SILENT_DATA);
      silentAudio.volume = 0;
      silentAudio.play().then(function() {
        silentAudio.pause();
        silentAudio = null; // メモリ解放
        
        audioUnlocked = true;
      }).catch(function() {});
    }
  }

  function playSound() {
    if (!chimeAudio) {
      initAudio();
      return;
    }

    try {
      chimeAudio.currentTime = 0;
      chimeAudio.volume = 1;
      chimeAudio.play().catch(function() {});
    } catch (e) {}
  }

  // フルスクリーン
  function toggleFullscreen() {
    var app = $('pokerTimer');
    if (!app) return;

    // タッチデバイス判定（iPad/iPhone/Android）
    var isTouchDevice = 'ontouchstart' in window || navigator.maxTouchPoints > 0;
    // iPad判定（iOS 13以降はMacとして報告されるため、タッチ+Macで判定）
    var isIPad = isTouchDevice ? (navigator.platform === 'MacIntel' || /iPad/.test(navigator.userAgent)) : false;
    var isIPhone = /iPhone/.test(navigator.userAgent);
    var isAndroid = /Android/.test(navigator.userAgent);
    var isMobileDevice = isIPad || isIPhone || isAndroid;

    // ネイティブFullscreen APIが使えるかチェック（デスクトップのみ使用）
    var canFullscreen = !isMobileDevice ? (app.requestFullscreen || app.webkitRequestFullscreen) : false;

    if (canFullscreen) {
      // デスクトップ: ネイティブFullscreen API
      var isFullscreen = document.fullscreenElement || document.webkitFullscreenElement;
      if (!isFullscreen) {
        if (app.requestFullscreen) {
          app.requestFullscreen();
        } else if (app.webkitRequestFullscreen) {
          app.webkitRequestFullscreen();
        }
      } else {
        if (document.exitFullscreen) {
          document.exitFullscreen();
        } else if (document.webkitExitFullscreen) {
          document.webkitExitFullscreen();
        }
      }
    } else {
      // モバイル: CSS擬似フルスクリーン
      var isFullscreen = app.classList.contains('mobile-fullscreen');
      if (!isFullscreen) {
        // フルスクリーンに入る
        app.classList.add('mobile-fullscreen');
        document.body.classList.add('mobile-fullscreen-active');
        // クラス追加後に高さを設定
        setMobileFullscreenHeight();
        window.addEventListener('resize', setMobileFullscreenHeight);
        window.addEventListener('orientationchange', setMobileFullscreenHeight);
      } else {
        // フルスクリーンを終了
        app.style.height = '';
        app.style.maxHeight = '';
        app.style.width = '';
        app.style.maxWidth = '';
        // timer-gridのスタイルもリセット
        var grid = app.querySelector('.timer-grid');
        if (grid) {
          grid.style.width = '';
          grid.style.maxWidth = '';
        }
        // 各パネルのスタイルもリセット
        var leftPanel = app.querySelector('.left-panel');
        var centerPanel = app.querySelector('.center-panel');
        var rightPanel = app.querySelector('.right-panel');
        if (leftPanel) {
          leftPanel.style.width = '';
          leftPanel.style.maxWidth = '';
        }
        if (centerPanel) {
          centerPanel.style.width = '';
          centerPanel.style.maxWidth = '';
        }
        if (rightPanel) {
          rightPanel.style.width = '';
          rightPanel.style.maxWidth = '';
        }
        app.classList.remove('mobile-fullscreen');
        document.body.classList.remove('mobile-fullscreen-active');
        window.removeEventListener('resize', setMobileFullscreenHeight);
        window.removeEventListener('orientationchange', setMobileFullscreenHeight);
      }
    }
  }

  // Safari/iOS対応: 実際のビューポートサイズを設定
  function setMobileFullscreenHeight() {
    var app = $('pokerTimer');
    if (!app) return;
    // window.innerHeight/innerWidthは実際の表示領域サイズを返す（Safari UIを除いた領域）
    var vh = window.innerHeight;
    var vw = window.innerWidth;
    app.style.height = vh + 'px';
    app.style.maxHeight = vh + 'px';
    app.style.width = vw + 'px';
    app.style.maxWidth = vw + 'px';

    // timer-gridの幅も明示的に設定（iPad対応）
    var grid = app.querySelector('.timer-grid');
    if (grid) {
      // パディングを考慮した幅を計算
      var paddingLeft = parseInt(window.getComputedStyle(app).paddingLeft) || 8;
      var paddingRight = parseInt(window.getComputedStyle(app).paddingRight) || 8;
      var gridWidth = vw - paddingLeft - paddingRight;
      grid.style.width = gridWidth + 'px';
      grid.style.maxWidth = gridWidth + 'px';
      // 各パネルの幅はCSSのFlexboxで制御（パーセンテージ指定）
    }
  }

  // キーボード操作
  function handleKeyDown(e) {
    if (e.target.tagName === 'INPUT') return;

    switch(e.key) {
      case 'ArrowLeft':
        adjustEntries(-1);
        e.preventDefault();
        break;
      case 'ArrowRight':
        adjustEntries(1);
        e.preventDefault();
        break;
      case 'ArrowUp':
        adjustRemaining(1);
        e.preventDefault();
        break;
      case 'ArrowDown':
        adjustRemaining(-1);
        e.preventDefault();
        break;
      case ' ':
        toggleTimer();
        e.preventDefault();
        break;
    }
  }

  // イベントリスナーを設定（ヘルパー関数）
  function addClick(id, handler) {
    var el = $(id);
    if (el) el.addEventListener('click', handler);
  }

  function setupEventListeners() {
    // メインコントロール
    addClick('startBtn', toggleTimer);
    addClick('btnPrev', prevLevel);
    addClick('btnSkip', skipLevel);
    addClick('btnSettings', openSettings);
    addClick('btnFullscreen', toggleFullscreen);
    addClick('btnFullscreenMobile', toggleFullscreen);
    addClick('btnExitFullscreen', toggleFullscreen);

    // プレイヤー調整
    addClick('btnEntryMinus', function() { adjustEntries(-1); });
    addClick('btnEntryPlus', function() { adjustEntries(1); });
    addClick('btnRemainMinus', function() { adjustRemaining(-1); });
    addClick('btnRemainPlus', function() { adjustRemaining(1); });

    // 設定モーダル
    addClick('btnSaveSettings', saveSettings);
    addClick('btnCloseSettings', closeSettings);
    addClick('btnSaveSet', saveBlindSet);
    addClick('btnDeleteSet', deleteBlindSet);
    addClick('btnAddLevel', addBlindLevel);
    addClick('btnAddBreak', addBreakLevel);
    addClick('btnCalcPrize', calculatePrizes);
    addClick('btnAddPrize', addPrizeRow);

    var blindSetSelect = $('blindSetSelect');
    if (blindSetSelect) blindSetSelect.addEventListener('change', loadBlindSet);

    // タブ切り替え
    var tabs = document.querySelectorAll('.settings-tab');
    for (var i = 0; i < tabs.length; i++) {
      tabs[i].addEventListener('click', function(e) {
        var tabName = e.target.getAttribute('data-tab');
        if (tabName) switchTab(tabName);
      });
    }

    // キーボード操作
    document.addEventListener('keydown', handleKeyDown);

    // iOS/iPad対応: 最初のタッチ/クリックでAudioContextを初期化
    var unlockAudio = function() {
      initAudio();
      document.removeEventListener('touchstart', unlockAudio);
      document.removeEventListener('touchend', unlockAudio);
      document.removeEventListener('click', unlockAudio);
    };
    document.addEventListener('touchstart', unlockAudio);
    document.addEventListener('touchend', unlockAudio);
    document.addEventListener('click', unlockAudio);
  }

  // 初期化
  function init() {
    setupEventListeners();
    initApp();
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
</script>
