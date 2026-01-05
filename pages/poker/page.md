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
/* メインアプリ */
.poker-timer-app {
  max-width: 1100px;
  margin: 0 auto;
  font-family: 'Inter', 'Helvetica Neue', Arial, sans-serif;
  background: #f8fafc;
  color: #1e293b;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.1);
  position: relative;
  touch-action: manipulation;  /* ダブルタップズームを無効化 */
}

.poker-timer-app * {
  box-sizing: border-box;
}

/* 3カラムレイアウト */
.timer-grid {
  display: grid;
  grid-template-columns: 220px 1fr 240px;
  gap: 15px;
  min-height: 400px;
}

/* 左カラム - PRIZE */
.left-panel {
  background: #fff;
  border-radius: 12px;
  padding: 15px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  display: flex;
  flex-direction: column;
}

.prize-header {
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 16px;
  margin-bottom: 16px;
}

.prize-inmoney {
  font-size: 18px;
  font-weight: 600;
  color: #475569;
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
  color: #475569;
  font-weight: 600;
}

.prize-amount {
  color: #22c55e;
  font-weight: 700;
}

/* 画像領域（プライズ下に固定表示） */
.qr-code-area {
  display: block;
  margin-top: auto;  /* 下に固定 */
  padding-top: 12px;
  text-align: center;
}

.qr-code-img {
  width: 180px;
  height: auto;
  border-radius: 12px;
  opacity: 0.95;
}

/* 中央カラム - タイマー */
.center-panel {
  background: linear-gradient(135deg, #1e293b 0%, #334155 100%);
  border-radius: 12px;
  padding: 25px;
  color: #fff;
  text-align: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
  position: relative;  /* フルスクリーンボタン配置用 */
}

/* レベル表示 - 大きく見やすく */
.level-badge {
  display: inline-block;
  background: #3b82f6;
  color: #fff;
  padding: 12px 32px;
  border-radius: 30px;
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 15px;
  letter-spacing: 1px;
}

.level-badge.break {
  background: #22c55e;
}

/* タイマー表示 - 大きく見やすく */
.timer-time {
  font-size: 96px;
  font-weight: 700;
  font-family: 'JetBrains Mono', 'Courier New', monospace;
  letter-spacing: 2px;
  margin: 20px 0;
  color: #ffffff;
  text-shadow: 0 4px 15px rgba(0,0,0,0.4);
}

.timer-time.warning {
  color: #ef4444;
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
  font-size: 48px;
  font-weight: 700;
  color: #fbbf24;
  margin: 16px 0;
  text-shadow: 0 3px 12px rgba(251, 191, 36, 0.4);
  display: flex;
  align-items: baseline;
  justify-content: center;
  gap: 20px;
  flex-wrap: wrap;
}

.blind-current .ante-value {
  font-size: 48px;
  font-weight: 700;
  color: #fbbf24;
  text-shadow: 0 3px 12px rgba(251, 191, 36, 0.4);
}

.blind-ante {
  display: none;
}

.blind-next {
  font-size: 18px;
  color: #94a3b8;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 2px solid rgba(255,255,255,0.2);
}

.blind-next-label {
  font-size: 16px;
  font-weight: 600;
  color: #cbd5e1;
  margin-bottom: 6px;
}

.blind-next-value {
  font-size: 28px;
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
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.btn:hover {
  transform: translateY(-2px);
}

.btn-primary {
  background: #22c55e;
  color: #fff;
}

.btn-primary:hover {
  background: #16a34a;
}

.btn-secondary {
  background: #64748b;
  color: #fff;
}

.btn-secondary:hover {
  background: #475569;
}

.btn-warning {
  background: #f59e0b;
  color: #fff;
}

.btn-warning:hover {
  background: #d97706;
}

/* 右カラム */
.right-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
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
  gap: 12px;
}

/* フルスクリーンボタン（PC用・right-panel右上固定） */
.fullscreen-btn-top {
  position: absolute;
  top: 0;
  right: 0;
  width: 36px;
  height: 36px;
  background: #3b82f6;
  color: #fff;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
  transition: background 0.2s;
}

.fullscreen-btn-top:hover {
  background: #2563eb;
}

/* フルスクリーンボタン（モバイル用・center-panelの右上に絶対配置） */
.fullscreen-btn-mobile {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 32px;
  height: 32px;
  background: rgba(59, 130, 246, 0.9);
  color: #fff;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  display: none;  /* PCでは非表示 */
  align-items: center;
  justify-content: center;
  z-index: 20;
  transition: background 0.2s;
}

.fullscreen-btn-mobile:hover {
  background: #2563eb;
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
  background: #fff;
  border-radius: 16px;
  padding: 20px;
  text-align: center;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
}

/* ラベルスタイル - 大きく見やすく */
.panel-label {
  font-size: 16px;
  font-weight: 700;
  color: #1e293b;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  margin-bottom: 10px;
}

.panel-value {
  font-size: 32px;
  font-weight: 700;
  color: #1e293b;
}

.panel-value.gold {
  color: #d97706;
}

.panel-value.large {
  font-size: 40px;
}

/* NEXT BREAK IN */
.break-card {
  background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
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
  background: linear-gradient(135deg, #64748b 0%, #475569 100%);
}

.break-card.no-break .panel-value {
  color: rgba(255,255,255,0.6);
}

/* STACK カード */
.stack-card {
  padding: 20px;
}

.stack-card .stack-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.stack-card .stack-label {
  font-size: 16px;
  font-weight: 500;
  color: #475569;
}

.stack-card .stack-value {
  font-size: 32px;
  font-weight: 700;
  color: #1e293b;
}

/* PLAYERS カード */
.players-card {
  padding: 20px;
}

.players-card .players-display {
  font-size: 32px;
  font-weight: 700;
  color: #1e293b;
  white-space: nowrap;
  margin: 12px 0;
}

.players-card .players-label {
  font-size: 14px;
  color: #475569;
  margin-bottom: 12px;
}

/* カウンターコントロール */
.counter-controls {
  display: flex;
  justify-content: center;
  gap: 8px;
}

.counter-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: 2px solid #e2e8f0;
  background: #fff;
  color: #475569;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.counter-btn:hover {
  background: #3b82f6;
  border-color: #3b82f6;
  color: #fff;
}

/* 設定モーダル - z-index最大 */
.settings-modal {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.6);
  z-index: 2147483647;
  overflow-y: auto;
}

.settings-modal.active {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 20px;
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
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  max-width: 650px;
  width: 100%;
  margin: 20px 0;
  color: #1e293b;
  max-height: 90vh;
  overflow-y: auto;
}

.settings-title {
  font-size: 22px;
  font-weight: bold;
  margin-bottom: 20px;
  text-align: center;
  color: #1e293b;
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
  color: #64748b;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
}

.settings-tab.active {
  color: #3b82f6;
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
  color: #475569;
  font-size: 14px;
}

.setting-input {
  width: 100%;
  padding: 10px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  background: #f8fafc;
  color: #1e293b;
  font-size: 16px;
}

.setting-input:focus {
  outline: none;
  border-color: #3b82f6;
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
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 16px;
}

.blind-set-selector button {
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
}

.btn-save-set {
  background: #22c55e;
  color: #fff;
}

.btn-delete-set {
  background: #ef4444;
  color: #fff;
}

/* ブラインドレベル一覧 */
.blind-levels {
  max-height: 280px;
  overflow-y: auto;
  margin-bottom: 12px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
}

.blind-level-item {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 8px;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
  cursor: grab;
  transition: background 0.2s, transform 0.2s;
  flex-wrap: wrap;
}

.blind-level-item:active {
  cursor: grabbing;
}

.blind-level-item.dragging {
  opacity: 0.5;
  background: #e2e8f0;
}

.blind-level-item.drag-over {
  border-top: 3px solid #3b82f6;
}

.blind-level-item:last-child {
  border-bottom: none;
}

.blind-level-item.break-item {
  background: #dcfce7;
}

.blind-level-item input {
  width: 96px;
  padding: 4px;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  background: #fff;
  color: #1e293b;
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
  width: 96px;
}

.blind-level-item .drag-handle {
  cursor: grab;
  color: #94a3b8;
  font-size: 14px;
  padding: 0 4px;
}

.blind-level-item .drag-handle:active {
  cursor: grabbing;
}

.blind-level-item .level-num {
  width: 20px;
  font-weight: bold;
  color: #3b82f6;
  font-size: 13px;
}

.blind-level-item .break-label {
  color: #22c55e;
  font-weight: 600;
  font-size: 14px;
}

.blind-level-item .delete-level {
  background: #ef4444;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 4px 8px;
  cursor: pointer;
  font-size: 14px;
  margin-left: auto;
}

.blind-level-item .time-label {
  font-size: 12px;
  color: #64748b;
}

.blind-level-item span {
  font-size: 12px;
  color: #64748b;
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
  background: #3b82f6;
  color: #fff;
}

.btn-add-break {
  background: #22c55e;
  color: #fff;
}

/* プライズ編集リスト */
.prize-edit-list {
  max-height: 250px;
  overflow-y: auto;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  margin-top: 10px;
}

.prize-edit-item {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 8px;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
  cursor: grab;
  transition: background 0.2s;
  flex-wrap: wrap;
}

.prize-edit-item:active {
  cursor: grabbing;
}

.prize-edit-item.dragging {
  opacity: 0.5;
  background: #e2e8f0;
}

.prize-edit-item.drag-over {
  border-top: 3px solid #3b82f6;
}

.prize-edit-item:last-child {
  border-bottom: none;
}

.prize-edit-item .drag-handle {
  cursor: grab;
  color: #94a3b8;
  font-size: 14px;
  padding: 0 4px;
}

.prize-edit-item .drag-handle:active {
  cursor: grabbing;
}

.prize-edit-item input {
  padding: 4px;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  font-size: 16px;
  box-sizing: border-box;
  text-align: center;
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
  color: #64748b;
}

.prize-edit-item .delete-prize {
  background: #ef4444;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 4px 8px;
  cursor: pointer;
  font-size: 14px;
  margin-left: auto;
}

.prize-edit-actions {
  display: flex;
  gap: 8px;
  margin-top: 10px;
}

.btn-add-prize {
  flex: 1;
  padding: 10px;
  background: #3b82f6;
  color: #fff;
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
  background: #3b82f6;
  color: #fff;
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

/* フルスクリーン時のスタイル */
.poker-timer-app:fullscreen {
  padding: 40px;
  display: flex;
  flex-direction: column;
  width: 100vw !important;
  height: 100vh !important;
  max-width: none !important;
  border-radius: 0 !important;
  background: #f8fafc !important;
}

.poker-timer-app:fullscreen .timer-grid {
  flex: 1;
  width: 100%;
  height: 100%;
  min-height: 0;
}

/* フルスクリーン時は設定ボタンのみ非表示（操作ボタンは表示） */
.poker-timer-app:fullscreen .controls #btnSettings {
  display: none !important;
}

/* フルスクリーン時のレイアウト調整 */
.poker-timer-app:fullscreen .timer-grid {
  gap: 16px;
  grid-template-columns: minmax(200px, 1.2fr) 2fr minmax(200px, 1.2fr);
  padding: 16px;
  box-sizing: border-box;
}

/* 中央パネル - タイマー */
.poker-timer-app:fullscreen .timer-time {
  font-size: 180px;
  font-weight: 700;
  color: #ffffff;
  text-shadow: 0 4px 20px rgba(0,0,0,0.3);
}

.poker-timer-app:fullscreen .timer-level {
  font-size: 36px;
  font-weight: 600;
  color: #ffffff;
}

.poker-timer-app:fullscreen .level-badge {
  font-size: 28px;
  padding: 14px 36px;
}

.poker-timer-app:fullscreen .blind-current {
  font-size: 72px;
  font-weight: 700;
  color: #fbbf24;
}

.poker-timer-app:fullscreen .blind-current .ante-value {
  font-size: 72px;
  color: #fbbf24;
}

.poker-timer-app:fullscreen .blind-next {
  margin-top: 16px;
}

.poker-timer-app:fullscreen .blind-next-label {
  font-size: 22px;
  color: #94a3b8;
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

/* 左パネル - PRIZE */
.poker-timer-app:fullscreen .left-panel {
  min-width: 0;
  height: 100%;
  overflow: hidden;
}

.poker-timer-app:fullscreen .left-panel {
  padding: 24px;
}

.poker-timer-app:fullscreen .left-panel .panel-label {
  font-size: 28px;
  font-weight: 600;
  color: #1e293b;
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
  color: #1e293b;
}

.poker-timer-app:fullscreen .prize-item .prize-amount {
  font-size: 24px;
  font-weight: 600;
  color: #22c55e;
}

/* 右パネル - NEXT BREAK */
.poker-timer-app:fullscreen .right-panel {
  min-width: 0;
  overflow: hidden;
}

.poker-timer-app:fullscreen .break-card {
  padding: 20px;
}

.poker-timer-app:fullscreen .break-card .panel-label {
  font-size: 28px;
  font-weight: 600;
  white-space: nowrap;
}

.poker-timer-app:fullscreen .break-card .panel-value {
  font-size: 36px;
  font-weight: 700;
  white-space: nowrap;
}

/* 右パネル - STACK */
.poker-timer-app:fullscreen .stack-card {
  padding: 20px;
}

.poker-timer-app:fullscreen .stack-card .panel-label {
  font-size: 28px;
  font-weight: 600;
  color: #1e293b;
}

.poker-timer-app:fullscreen .stack-row .stack-label {
  font-size: 24px;
  color: #475569;
}

.poker-timer-app:fullscreen .stack-row .stack-value {
  font-size: 48px;
  font-weight: 700;
  color: #1e293b;
}

/* 右パネル - PLAYERS */
.poker-timer-app:fullscreen .players-card {
  padding: 20px;
}

.poker-timer-app:fullscreen .players-card .panel-label {
  font-size: 28px;
  font-weight: 600;
  color: #1e293b;
}

.poker-timer-app:fullscreen .players-display {
  font-size: 48px;
  font-weight: 700;
  color: #1e293b;
  white-space: nowrap;
}

.poker-timer-app:fullscreen .players-label {
  font-size: 22px;
  color: #475569;
}

.poker-timer-app:fullscreen .players-buttons {
  display: none;
}

/* フルスクリーンボタン非表示 */
.poker-timer-app:fullscreen .fullscreen-btn-top {
  display: none;
}

/* フルスクリーン時の画像拡大 */
.poker-timer-app:fullscreen .qr-code-img {
  width: 280px;
  height: auto;
}

/* モバイル擬似フルスクリーン */
body.mobile-fullscreen-active {
  overflow: hidden !important;
}

.poker-timer-app.mobile-fullscreen {
  position: fixed !important;
  top: 0 !important;
  left: 0 !important;
  width: 100vw !important;
  /* heightはJavaScriptでwindow.innerHeightを設定 */
  max-width: none !important;
  border-radius: 0 !important;
  z-index: 2147483647 !important;
  padding: 8px;
  padding-top: max(8px, env(safe-area-inset-top));
  padding-bottom: max(8px, env(safe-area-inset-bottom));
  padding-left: max(8px, env(safe-area-inset-left));
  padding-right: max(8px, env(safe-area-inset-right));
  display: flex;
  flex-direction: column;
  background: #f8fafc !important;
  box-sizing: border-box;
  overflow: hidden;
}

.poker-timer-app.mobile-fullscreen .timer-grid {
  flex: 1;
  width: 100%;
  min-height: 0;
  display: grid;
  /* デフォルトは1カラム（iPhone向け） */
  grid-template-columns: 1fr;
  grid-template-rows: auto 1fr auto;
  padding: 4px;
  gap: 4px;
  box-sizing: border-box;
  align-items: stretch;
}

/* iPad以上（横幅600px以上）では3カラム */
@media (min-width: 600px) {
  .poker-timer-app.mobile-fullscreen .timer-grid {
    grid-template-columns: minmax(80px, 1fr) 2fr minmax(80px, 1fr);
    grid-template-rows: 1fr;
  }
}

/* モバイルフルスクリーン時は設定ボタンのみ非表示 */
.poker-timer-app.mobile-fullscreen .controls #btnSettings {
  display: none !important;
}

/* モバイルフルスクリーン時のフォントサイズ */

/* 中央パネル */
.poker-timer-app.mobile-fullscreen .center-panel {
  min-height: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;  /* 縦方向中央揃え */
}

.poker-timer-app.mobile-fullscreen .timer-card {
  padding: 6px;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 0;
}

/* 中央パネル - タイマー */
.poker-timer-app.mobile-fullscreen .timer-time {
  font-size: 100px;
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
  font-size: 20px;
}

.poker-timer-app.mobile-fullscreen .blind-current {
  font-size: 40px;
  font-weight: 700;
  color: #fbbf24;
}

.poker-timer-app.mobile-fullscreen .blind-current .ante-value {
  font-size: 40px;
  color: #fbbf24;
}

.poker-timer-app.mobile-fullscreen .blind-next-label {
  font-size: 16px;
}

.poker-timer-app.mobile-fullscreen .blind-next-value {
  font-size: 24px;
  font-weight: 600;
  color: #e2e8f0;
}

.poker-timer-app.mobile-fullscreen .next-ante {
  font-size: 16px;
  color: #e2e8f0;
}

.poker-timer-app.mobile-fullscreen .timer-progress {
  height: 6px;
  margin: 4px 0;
}

/* 左パネル - PRIZE */
.poker-timer-app.mobile-fullscreen .left-panel {
  min-width: 0;
  min-height: 0;
  height: 100%;
}

.poker-timer-app.mobile-fullscreen .left-panel .panel-label {
  font-size: 22px;
  margin-bottom: 2px;
}

.poker-timer-app.mobile-fullscreen .prize-inmoney {
  font-size: 22px;
}

.poker-timer-app.mobile-fullscreen .prize-inmoney-value {
  font-size: 28px;
}

.poker-timer-app.mobile-fullscreen .prize-item {
  font-size: 24px;
  padding: 3px 0;
}

/* 右パネル */
.poker-timer-app.mobile-fullscreen .right-panel {
  min-width: 0;
  min-height: 0;
  gap: 2px;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;  /* 上から順に詰める */
}

/* 右パネル - NEXT BREAK */
.poker-timer-app.mobile-fullscreen .break-card {
  padding: 4px 6px;
  flex: 0 1 auto;
  min-height: 0;
  overflow: hidden;
}

.poker-timer-app.mobile-fullscreen .break-card .panel-label {
  font-size: 22px;
  white-space: nowrap;
  margin-bottom: 2px;
}

.poker-timer-app.mobile-fullscreen .break-card .panel-value {
  font-size: 36px;
  white-space: nowrap;
}

/* 右パネル - スペーサー: モバイルフルスクリーンでは非表示 */
.poker-timer-app.mobile-fullscreen .right-panel-spacer {
  display: none !important;
}

/* 右パネル - 下部エリア */
.poker-timer-app.mobile-fullscreen .right-panel-bottom {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 0 0 auto;  /* 固定サイズ、縮小しない */
  min-height: 0;
}

/* 右パネル - STACK */
.poker-timer-app.mobile-fullscreen .stack-card {
  padding: 4px 6px;
  flex: 0 1 auto;
  min-height: 0;
  overflow: hidden;
}

.poker-timer-app.mobile-fullscreen .stack-card .panel-label {
  font-size: 22px;
  margin-bottom: 2px;
}

.poker-timer-app.mobile-fullscreen .stack-row .stack-label {
  font-size: 22px;
}

.poker-timer-app.mobile-fullscreen .stack-row .stack-value {
  font-size: 36px;
}

.poker-timer-app.mobile-fullscreen .stack-row {
  margin: 1px 0;
}

/* 右パネル - PLAYERS */
.poker-timer-app.mobile-fullscreen .players-card {
  padding: 4px 6px;
  flex: 0 1 auto;
  min-height: 0;
  overflow: hidden;
}

.poker-timer-app.mobile-fullscreen .players-card .panel-label {
  font-size: 22px;
  margin-bottom: 2px;
}

.poker-timer-app.mobile-fullscreen .players-display {
  font-size: 44px;
  font-weight: 700;
  white-space: nowrap;
}

.poker-timer-app.mobile-fullscreen .players-label {
  font-size: 20px;
  color: #475569;
}

.poker-timer-app.mobile-fullscreen .players-buttons {
  display: none;
}

.poker-timer-app.mobile-fullscreen .fullscreen-btn-top {
  display: none;
}

.poker-timer-app.mobile-fullscreen .fullscreen-btn-mobile {
  display: none;
}

/* フルスクリーン解除ボタン（モバイルフルスクリーン時のみ表示） */
.exit-fullscreen-btn {
  display: none;
}

.poker-timer-app.mobile-fullscreen .exit-fullscreen-btn {
  display: flex;
  position: fixed;
  top: max(8px, env(safe-area-inset-top));
  right: max(8px, env(safe-area-inset-right));
  width: 40px;
  height: 40px;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 18px;
  cursor: pointer;
  z-index: 2147483647;
  align-items: center;
  justify-content: center;
}

.poker-timer-app.mobile-fullscreen .exit-fullscreen-btn:hover {
  background: rgba(0, 0, 0, 0.8);
}

/* スマートフォン向けモバイルフルスクリーン（縦画面600px未満） */
@media (max-width: 599px) and (orientation: portrait) {
  .poker-timer-app.mobile-fullscreen .left-panel {
    display: none !important;  /* PRIZEを非表示 */
  }

  .poker-timer-app.mobile-fullscreen .center-panel {
    order: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    min-height: 0;
    flex: 1;
  }

  .poker-timer-app.mobile-fullscreen .right-panel {
    order: 2;
    display: grid;
    grid-template-columns: 1fr 1fr;  /* 2カラム */
    grid-template-rows: auto auto;
    gap: 6px;
    padding: 4px;
  }

  .poker-timer-app.mobile-fullscreen .right-panel-spacer {
    display: none !important;
  }

  .poker-timer-app.mobile-fullscreen .right-panel-bottom {
    display: contents;  /* 子要素を直接配置 */
  }

  /* NEXT BREAK: 左上 */
  .poker-timer-app.mobile-fullscreen .break-card {
    grid-column: 1;
    grid-row: 1;
    padding: 6px 10px;
  }

  /* STACK: 左下 */
  .poker-timer-app.mobile-fullscreen .stack-card {
    grid-column: 1;
    grid-row: 2;
    padding: 6px 10px;
  }

  /* PLAYERS: 右側（2行分） */
  .poker-timer-app.mobile-fullscreen .players-card {
    grid-column: 2;
    grid-row: 1 / 3;
    padding: 6px 10px;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  /* iPhone: タイマーを大きく */
  .poker-timer-app.mobile-fullscreen .timer-time {
    font-size: 64px;
  }

  .poker-timer-app.mobile-fullscreen .blind-current {
    font-size: 26px;
  }

  .poker-timer-app.mobile-fullscreen .blind-current .ante-value {
    font-size: 26px;
  }

  .poker-timer-app.mobile-fullscreen .blind-next-value {
    font-size: 16px;
  }

  .poker-timer-app.mobile-fullscreen .level-badge {
    font-size: 14px;
    padding: 3px 10px;
  }

  /* パネルラベルのフォント統一 */
  .poker-timer-app.mobile-fullscreen .left-panel .panel-label,
  .poker-timer-app.mobile-fullscreen .break-card .panel-label,
  .poker-timer-app.mobile-fullscreen .stack-card .panel-label,
  .poker-timer-app.mobile-fullscreen .players-card .panel-label {
    font-size: 12px;
  }

  .poker-timer-app.mobile-fullscreen .break-card .panel-value {
    font-size: 20px;
  }

  .poker-timer-app.mobile-fullscreen .stack-row .stack-label {
    font-size: 12px;
  }

  .poker-timer-app.mobile-fullscreen .stack-row .stack-value {
    font-size: 18px;
  }

  .poker-timer-app.mobile-fullscreen .players-display {
    font-size: 24px;
  }

  .poker-timer-app.mobile-fullscreen .players-label {
    font-size: 11px;
  }
}

/* スマートフォン向けモバイルフルスクリーン（横画面） */
@media (max-height: 500px) and (orientation: landscape) {
  .poker-timer-app.mobile-fullscreen {
    /* safe-areaを確保（iPhoneのノッチ・ホームインジケーター対応） */
    padding: 4px;
    padding-left: max(4px, env(safe-area-inset-left, 4px));
    padding-right: max(4px, env(safe-area-inset-right, 4px));
    padding-bottom: max(4px, env(safe-area-inset-bottom, 4px));
  }

  .poker-timer-app.mobile-fullscreen .timer-grid {
    grid-template-columns: minmax(100px, 1fr) 2fr minmax(100px, 1fr);  /* 3カラム: PRIZE, タイマー, 右パネル */
    grid-template-rows: 1fr;
    gap: 4px;
    height: 100%;
    max-height: 100%;
    overflow: hidden;
  }

  /* PRIZE（左カラム）を表示 */
  .poker-timer-app.mobile-fullscreen .left-panel {
    display: flex !important;
    flex-direction: column;
    order: 1;
    min-width: 0;
    max-height: 100%;
    overflow: hidden;
    padding: 4px;
  }

  .poker-timer-app.mobile-fullscreen .prize-card {
    padding: 2px 4px;
  }

  .poker-timer-app.mobile-fullscreen .prize-header {
    padding-bottom: 2px;
    margin-bottom: 2px;
  }

  .poker-timer-app.mobile-fullscreen .prize-inmoney {
    font-size: 10px;
  }

  .poker-timer-app.mobile-fullscreen .prize-inmoney-value {
    font-size: 12px;
  }

  .poker-timer-app.mobile-fullscreen .prize-list {
    /* max-height削除: 画像の上まで自動的に伸びる */
    flex: 1;
    overflow: hidden;
  }

  .poker-timer-app.mobile-fullscreen .prize-list-inner {
    position: relative;
    will-change: transform;
  }

  .poker-timer-app.mobile-fullscreen .prize-item {
    font-size: 11px;
    padding: 1px 0;
  }

  /* 横画面フルスクリーン時にQRコードを表示 */
  .poker-timer-app.mobile-fullscreen .qr-code-area {
    display: block;
    margin-top: auto;
    padding-top: 4px;
  }

  .poker-timer-app.mobile-fullscreen .qr-code-img {
    width: 100px;
    height: auto;
  }

  .poker-timer-app.mobile-fullscreen .center-panel {
    order: 2;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    min-height: 0;
    max-height: 100%;
    overflow: hidden;
    padding: 4px;
  }

  .poker-timer-app.mobile-fullscreen .right-panel {
    order: 3;
    display: flex;
    flex-direction: column;
    gap: 2px;
    padding: 4px;
    justify-content: flex-end;  /* 下揃えに変更 */
    min-width: 0;
    max-height: 100%;
    overflow: hidden;
  }

  .poker-timer-app.mobile-fullscreen .right-panel-spacer {
    display: none !important;
  }

  .poker-timer-app.mobile-fullscreen .right-panel-bottom {
    display: contents;
  }

  .poker-timer-app.mobile-fullscreen .break-card,
  .poker-timer-app.mobile-fullscreen .stack-card,
  .poker-timer-app.mobile-fullscreen .players-card {
    padding: 2px 4px;
  }

  /* 横画面: タイマーサイズ調整 */
  .poker-timer-app.mobile-fullscreen .timer-time {
    font-size: 56px;
    margin: 4px 0;
  }

  .poker-timer-app.mobile-fullscreen .blind-current {
    font-size: 22px;
    margin: 4px 0;
  }

  .poker-timer-app.mobile-fullscreen .blind-current .ante-value {
    font-size: 22px;
  }

  .poker-timer-app.mobile-fullscreen .blind-next-value {
    font-size: 14px;
  }

  .poker-timer-app.mobile-fullscreen .level-badge {
    font-size: 12px;
    padding: 3px 8px;
    margin-bottom: 4px;
  }

  .poker-timer-app.mobile-fullscreen .timer-progress {
    height: 3px;
    margin: 2px 0;
  }

  .poker-timer-app.mobile-fullscreen .blind-info {
    margin: 4px 0;
  }

  .poker-timer-app.mobile-fullscreen .blind-next {
    margin-top: 4px;
    padding-top: 4px;
  }

  /* パネルラベルのフォント統一 */
  .poker-timer-app.mobile-fullscreen .left-panel .panel-label,
  .poker-timer-app.mobile-fullscreen .break-card .panel-label,
  .poker-timer-app.mobile-fullscreen .stack-card .panel-label,
  .poker-timer-app.mobile-fullscreen .players-card .panel-label {
    font-size: 10px;
    margin-bottom: 1px;
  }

  .poker-timer-app.mobile-fullscreen .break-card .panel-value {
    font-size: 16px;
  }

  .poker-timer-app.mobile-fullscreen .stack-row .stack-label {
    font-size: 10px;
  }

  .poker-timer-app.mobile-fullscreen .stack-row .stack-value {
    font-size: 14px;
  }

  .poker-timer-app.mobile-fullscreen .stack-row {
    padding: 1px 0;
  }

  .poker-timer-app.mobile-fullscreen .players-display {
    font-size: 20px;
  }

  .poker-timer-app.mobile-fullscreen .players-label {
    font-size: 9px;
  }
}

/* 使い方説明 */
.usage-guide {
  max-width: 1100px;
  margin: 30px auto 0;
  padding: 24px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.08);
  color: #1e293b;
  font-size: 14px;
  line-height: 1.8;
}

.usage-guide h3 {
  margin: 0 0 16px 0;
  font-size: 18px;
  color: #1e293b;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 8px;
}

.usage-guide h4 {
  margin: 20px 0 8px 0;
  font-size: 15px;
  color: #475569;
}

.usage-guide ul {
  margin: 8px 0;
  padding-left: 20px;
}

.usage-guide li {
  margin: 4px 0;
}

.usage-guide kbd {
  display: inline-block;
  padding: 2px 6px;
  font-size: 12px;
  font-family: monospace;
  background: #f1f5f9;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
}

/* レスポンシブ */
@media (max-width: 900px) {
  .timer-grid {
    grid-template-columns: 1fr;
    grid-template-rows: auto auto auto;
    gap: 12px;
  }

  .center-panel {
    order: 1;
    padding: 12px;
  }

  /* モバイル: Levelバッジとフルスクリーンボタンを横並び表示 */
  .level-row {
    display: flex;
    margin-bottom: 4px;
  }

  /* モバイル: 単独のLevelバッジを非表示 */
  .center-panel > .level-badge {
    display: none;
  }

  .timer-time {
    margin: 6px 0;
  }

  .blind-info {
    margin-top: 6px;
  }

  .controls {
    margin-top: 10px;
  }

  /* モバイル: PC用ボタンを非表示、モバイル用ボタンを表示 */
  .fullscreen-btn-top {
    display: none;
  }

  .fullscreen-btn-mobile {
    display: flex;
  }

  /* 右パネル: 3行グリッド（NEXT BREAK上、STACK+PLAYERS下横並び） */
  .right-panel {
    order: 2;
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: auto auto;
    gap: 8px;
    padding-top: 0;
    height: auto;
  }

  .right-panel-spacer {
    display: none;
  }

  .right-panel-bottom {
    display: contents;
  }

  /* NEXT BREAK: 左上（フルスクリーン時と同じレイアウト） */
  .break-card {
    grid-column: 1;
    grid-row: 1;
  }

  /* STACK: 左下（フルスクリーン時と同じレイアウト） */
  .stack-card {
    grid-column: 1;
    grid-row: 2;
  }

  /* PLAYERS: 右側2行分（フルスクリーン時と同じレイアウト） */
  .players-card {
    grid-column: 2;
    grid-row: 1 / 3;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  /* PRIZEは一番下 */
  .left-panel {
    order: 3;
    height: 210px;
    max-height: 210px;
    overflow: hidden;
  }

  .level-badge {
    font-size: 14px;
    padding: 6px 12px;
    margin-top: 0;
    margin-bottom: 0;
  }

  .timer-time {
    font-size: 64px;
  }

  .blind-current {
    font-size: 32px;
  }

  .blind-current .ante-value {
    font-size: 32px;
  }

  .blind-next-value {
    font-size: 20px;
  }

  .panel-value {
    font-size: 24px;
  }

  .break-card .panel-value {
    font-size: 28px;
  }

  .players-card .players-display {
    font-size: 20px;
  }

  .prize-list {
    /* max-height削除: 画像の上まで自動的に伸びる */
    flex: 0 0 auto;
    height: 120px;
    overflow: hidden;
    position: relative;
    min-height: 40px;
  }

  /* モバイル時は画像を非表示 */
  .qr-code-area {
    display: none;
  }

  .prize-header {
    padding-bottom: 8px;
    margin-bottom: 8px;
  }

  .prize-list-inner {
    position: absolute;
    width: 100%;
    top: 0;
    left: 0;
  }

  .prize-item {
    padding: 4px 0;
    font-size: 16px;
  }

  .prize-inmoney {
    font-size: 16px;
  }

  /* info-cardをコンパクトに */
  .right-panel .info-card {
    padding: 10px 12px;
  }

  .right-panel .panel-label {
    font-size: 16px;
    margin-bottom: 4px;
  }

  .left-panel .panel-label {
    font-size: 16px;
  }

  /* STACKカードのフォントサイズ調整 */
  .stack-card .stack-row {
    padding: 4px 0;
  }

  .stack-card .stack-label {
    font-size: 12px;
  }

  .stack-card .stack-value {
    font-size: 20px;
  }
}

@media (max-width: 500px) {
  .poker-timer-app {
    padding: 10px;
  }

  .timer-grid {
    gap: 8px;
  }

  .timer-time {
    font-size: 56px;
  }

  .blind-current {
    font-size: 28px;
  }

  .blind-current .ante-value {
    font-size: 28px;
  }

  .level-badge {
    font-size: 16px;
    padding: 6px 16px;
  }

  /* 右パネルをさらにコンパクトに */
  .right-panel {
    gap: 6px;
  }

  .right-panel .info-card {
    padding: 8px 10px;
  }

  .right-panel .panel-label {
    font-size: 14px;
    margin-bottom: 2px;
  }

  .left-panel .panel-label {
    font-size: 14px;
  }

  .panel-value {
    font-size: 20px;
  }

  .break-card .panel-value {
    font-size: 24px;
  }

  .players-card .players-display {
    font-size: 16px;
  }

  /* STACKカードをさらにコンパクトに */
  .stack-card .stack-row {
    padding: 2px 0;
  }

  .stack-card .stack-label {
    font-size: 10px;
  }

  .stack-card .stack-value {
    font-size: 16px;
  }

  .players-buttons .btn {
    width: 32px;
    height: 32px;
    font-size: 14px;
  }

  /* PRIZEもコンパクトに */
  .left-panel {
    height: 210px;
    max-height: 210px;
    overflow: hidden;
  }

  .prize-item {
    font-size: 14px;
  }

  .prize-inmoney {
    font-size: 14px;
  }

  .prize-list {
    /* max-height削除: 画像の上まで自動的に伸びる */
    flex: 0 0 auto;
    height: 120px;
    overflow: hidden;
  }

  .prize-list-inner {
    position: absolute;
    width: 100%;
    top: 0;
    left: 0;
    will-change: transform;
  }

  .controls {
    flex-wrap: wrap;
    gap: 6px;
  }

  .btn {
    padding: 10px 14px;
    font-size: 13px;
  }
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
      <div class="qr-code-area" id="qrCodeArea">
        <img src="data:image/jpeg;base64,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" alt="ポーカーやろうよ" class="qr-code-img">
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
        <label class="setting-label">ブラインドレベル (ドラッグで並べ替え可能)</label>
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
        <label class="setting-label">プライズ配分 (ドラッグで並べ替え可能)</label>
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
    var currentBlind = state.blinds[state.currentLevel];
    if (currentBlind) {
      if (currentBlind.isBreak) {
        state.levelTime = (currentBlind.breakTime || 5) * 60;
      } else {
        state.levelTime = (currentBlind.levelTime || state.defaultLevelTime) * 60;
      }
      state.timeRemaining = state.levelTime;
    }

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
    var currentBlind = state.blinds[0];
    if (currentBlind) {
      if (currentBlind.isBreak) {
        state.levelTime = (currentBlind.breakTime || 5) * 60;
      } else {
        state.levelTime = (currentBlind.levelTime || defaultLevelTime) * 60;
      }
      state.timeRemaining = state.levelTime;
    }

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
    if (state.isRunning) {
      clearInterval(state.intervalId);
      state.isRunning = false;
      var startBtn = $('startBtn');
      if (startBtn) startBtn.textContent = '開始';
    } else {
      state.isRunning = true;
      var startBtn2 = $('startBtn');
      if (startBtn2) startBtn2.textContent = '停止';
      state.intervalId = setInterval(tick, 1000);
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
      var blind = state.blinds[state.currentLevel];
      if (blind.isBreak) {
        state.levelTime = (blind.breakTime || 5) * 60;
      } else {
        state.levelTime = (blind.levelTime || state.defaultLevelTime) * 60;
      }
      state.timeRemaining = state.levelTime;
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
      var blind = state.blinds[state.currentLevel];
      if (blind.isBreak) {
        state.levelTime = (blind.breakTime || 5) * 60;
      } else {
        state.levelTime = (blind.levelTime || state.defaultLevelTime) * 60;
      }
      state.timeRemaining = state.levelTime;
      updateDisplay();
    }
  }

  function skipLevel() {
    nextLevel();
  }

  // 表示更新
  function updateDisplay() {
    var minutes = Math.floor(state.timeRemaining / 60);
    var seconds = state.timeRemaining % 60;
    var td = $('timerDisplay');
    if (td) td.textContent = String(minutes).padStart(2, '0') + ':' + String(seconds).padStart(2, '0');

    // Warning表示
    if (state.timeRemaining <= 60) {
      if (td) td.classList.add('warning');
    } else {
      if (td) td.classList.remove('warning');
    }

    // プログレスバー
    var progress = state.levelTime > 0 ? (state.timeRemaining / state.levelTime) * 100 : 0;
    var pb = $('progressBar');
    if (pb) pb.style.width = progress + '%';

    // ブラインド情報
    var blind = state.blinds[state.currentLevel];
    var lb = $('levelBadge');
    var lbMobile = $('levelBadgeMobile');
    var cb = $('currentBlind');
    var ai = $('anteInfo');
    var nb = $('nextBlind');

    if (blind) {
      if (blind.isBreak) {
        if (lb) {
          lb.textContent = 'BREAK';
          lb.classList.add('break');
        }
        if (lbMobile) {
          lbMobile.textContent = 'BREAK';
          lbMobile.classList.add('break');
        }
        if (cb) cb.textContent = 'BREAK TIME';
        if (ai) ai.textContent = '';
      } else {
        var levelNum = 0;
        for (var i = 0; i <= state.currentLevel; i++) {
          if (!state.blinds[i].isBreak) levelNum++;
        }
        var totalLevelsCount = 0;
        for (var j = 0; j < state.blinds.length; j++) {
          if (!state.blinds[j].isBreak) totalLevelsCount++;
        }
        if (lb) {
          lb.innerHTML = 'Level <span id="currentLevel">' + levelNum + '</span> / <span id="totalLevels">' + totalLevelsCount + '</span>';
          lb.classList.remove('break');
        }
        if (lbMobile) {
          lbMobile.innerHTML = 'Level <span id="currentLevelMobile">' + levelNum + '</span> / <span id="totalLevelsMobile">' + totalLevelsCount + '</span>';
          lbMobile.classList.remove('break');
        }
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
    if (nextBlindData) {
      if (nb) {
        nb.innerHTML = nextBlindData.sb + ' / ' + nextBlindData.bb + ' <span class="next-ante">(Ante: ' + nextBlindData.ante + ')</span>';
      }
    } else {
      if (nb) nb.innerHTML = '--';
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
    var showBreak = false;
    if (foundBreak) {
      if (!state.blinds[state.currentLevel].isBreak) {
        showBreak = true;
      }
    }
    if (showBreak) {
      breakCard.classList.remove('no-break');
      var minutes = Math.floor(timeToBreak / 60);
      var seconds = timeToBreak % 60;
      if (nbt) nbt.textContent = String(minutes).padStart(2, '0') + ':' + String(seconds).padStart(2, '0');
    } else {
      // ブレイクがなくても常時表示（--:-- を表示）
      breakCard.classList.add('no-break');
      if (nbt) nbt.textContent = '--:--';
    }
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
            rankText = p.startRank + '〜' + p.endRank + '位';
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
      if (inner.getAttribute('data-duplicated') && inner.getAttribute('data-original-html')) {
        inner.innerHTML = inner.getAttribute('data-original-html');
        inner.removeAttribute('data-duplicated');
        inner.removeAttribute('data-original-html');
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

    var containerHeight = container.clientHeight;
    var innerHeight = inner.scrollHeight;
    if (!containerHeight || !innerHeight || innerHeight <= containerHeight + 1) return;

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
        '<span class="drag-handle">☰</span>' +
        '<input type="number" class="rank-input rank-start" value="' + p.startRank + '" min="1">' +
        '<span>位 〜</span>' +
        '<input type="number" class="rank-input rank-end" value="' + p.endRank + '" min="1">' +
        '<span>位:</span>' +
        '<input type="number" class="amount-input" value="' + p.amount + '" min="0" step="100">' +
        '<span>pt</span>' +
        '<button class="delete-prize">×</button>';
      container.appendChild(item);
    }

    setupPrizeDragAndDrop();
    setupPrizeDeleteButtons();
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

      // インデックスを再設定
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

    var newIndex = items.length;
    var item = document.createElement('div');
    item.className = 'prize-edit-item';
    item.setAttribute('data-index', newIndex);
    item.setAttribute('draggable', 'true');
    item.innerHTML =
      '<span class="drag-handle">☰</span>' +
      '<input type="number" class="rank-input rank-start" value="' + nextRank + '" min="1">' +
      '<span>位 〜</span>' +
      '<input type="number" class="rank-input rank-end" value="' + nextRank + '" min="1">' +
      '<span>位:</span>' +
      '<input type="number" class="amount-input" value="0" min="0" step="100">' +
      '<span>pt</span>' +
      '<button class="delete-prize">×</button>';

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
    if (items.length > 1 && items[index]) {
      items[index].remove();
      // インデックスを再設定
      var newItems = container.querySelectorAll('.prize-edit-item');
      for (var i = 0; i < newItems.length; i++) {
        newItems[i].setAttribute('data-index', i);
      }
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
          '<span class="drag-handle">☰</span>' +
          '<span class="break-label">BREAK</span>' +
          '<input type="number" class="break-time-input" value="' + (blind.breakTime || 5) + '" min="1" style="width: 50px;"> 分' +
          '<button class="delete-level">×</button>';
      } else {
        levelNum++;
        item.className = 'blind-level-item';
        item.setAttribute('data-index', index);
        item.setAttribute('draggable', 'true');
        item.innerHTML =
          '<span class="drag-handle">☰</span>' +
          '<span class="level-num">' + levelNum + '</span>' +
          '<input type="number" class="time-input" value="' + (blind.levelTime || state.defaultLevelTime) + '" min="1" title="時間(分)">' +
          '<span class="time-label">分</span>' +
          '<input type="number" class="sb-input" value="' + blind.sb + '" placeholder="SB">' +
          '<span>/</span>' +
          '<input type="number" class="bb-input" value="' + blind.bb + '" placeholder="BB">' +
          '<span>/</span>' +
          '<input type="number" class="ante-input" value="' + blind.ante + '" placeholder="Ante">' +
          '<button class="delete-level">×</button>';
      }
      container.appendChild(item);
    }

    setupBlindDragAndDrop();
    setupBlindDeleteButtons();
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

      // state.blindsの配列も並べ替え
      var movedBlind = state.blinds.splice(fromIndex, 1)[0];
      state.blinds.splice(toIndex, 0, movedBlind);

      // 画面を再描画
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

  function playSound() {
    try {
      var audioCtx = new (window.AudioContext || window.webkitAudioContext)();
      var oscillator = audioCtx.createOscillator();
      var gainNode = audioCtx.createGain();
      oscillator.connect(gainNode);
      gainNode.connect(audioCtx.destination);
      oscillator.frequency.value = 800;
      oscillator.type = 'sine';
      gainNode.gain.value = 0.3;
      oscillator.start();
      setTimeout(function() { oscillator.stop(); }, 200);
    } catch (e) {}
  }

  // フルスクリーン
  function toggleFullscreen() {
    var app = $('pokerTimer');
    if (!app) return;

    // ネイティブFullscreen APIが使えるかチェック
    var canFullscreen = app.requestFullscreen || app.webkitRequestFullscreen;

    if (canFullscreen) {
      // デスクトップ: ネイティブFullscreen API
      if (!document.fullscreenElement && !document.webkitFullscreenElement) {
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
        app.classList.remove('mobile-fullscreen');
        document.body.classList.remove('mobile-fullscreen-active');
        window.removeEventListener('resize', setMobileFullscreenHeight);
        window.removeEventListener('orientationchange', setMobileFullscreenHeight);
      }
    }
  }

  // Safari/iOS対応: 実際のビューポート高さを設定
  function setMobileFullscreenHeight() {
    var app = $('pokerTimer');
    if (!app) return;
    // window.innerHeightは実際の表示領域高さを返す（Safari UIを除いた領域）
    var vh = window.innerHeight;
    app.style.height = vh + 'px';
    app.style.maxHeight = vh + 'px';
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

  // イベントリスナーを設定
  function setupEventListeners() {
    // メインコントロール
    var startBtn = $('startBtn');
    if (startBtn) startBtn.addEventListener('click', toggleTimer);

    var btnPrev = $('btnPrev');
    if (btnPrev) btnPrev.addEventListener('click', prevLevel);

    var btnSkip = $('btnSkip');
    if (btnSkip) btnSkip.addEventListener('click', skipLevel);

    var btnSettings = $('btnSettings');
    if (btnSettings) btnSettings.addEventListener('click', openSettings);

    var btnFullscreen = $('btnFullscreen');
    if (btnFullscreen) btnFullscreen.addEventListener('click', toggleFullscreen);

    var btnFullscreenMobile = $('btnFullscreenMobile');
    if (btnFullscreenMobile) btnFullscreenMobile.addEventListener('click', toggleFullscreen);

    var btnExitFullscreen = $('btnExitFullscreen');
    if (btnExitFullscreen) btnExitFullscreen.addEventListener('click', toggleFullscreen);

    // プレイヤー調整
    var btnEntryMinus = $('btnEntryMinus');
    if (btnEntryMinus) btnEntryMinus.addEventListener('click', function() { adjustEntries(-1); });

    var btnEntryPlus = $('btnEntryPlus');
    if (btnEntryPlus) btnEntryPlus.addEventListener('click', function() { adjustEntries(1); });

    var btnRemainMinus = $('btnRemainMinus');
    if (btnRemainMinus) btnRemainMinus.addEventListener('click', function() { adjustRemaining(-1); });

    var btnRemainPlus = $('btnRemainPlus');
    if (btnRemainPlus) btnRemainPlus.addEventListener('click', function() { adjustRemaining(1); });

    // 設定モーダル
    var btnSaveSettings = $('btnSaveSettings');
    if (btnSaveSettings) btnSaveSettings.addEventListener('click', saveSettings);

    var btnCloseSettings = $('btnCloseSettings');
    if (btnCloseSettings) btnCloseSettings.addEventListener('click', closeSettings);

    var btnSaveSet = $('btnSaveSet');
    if (btnSaveSet) btnSaveSet.addEventListener('click', saveBlindSet);

    var btnDeleteSet = $('btnDeleteSet');
    if (btnDeleteSet) btnDeleteSet.addEventListener('click', deleteBlindSet);

    var blindSetSelect = $('blindSetSelect');
    if (blindSetSelect) blindSetSelect.addEventListener('change', loadBlindSet);

    var btnAddLevel = $('btnAddLevel');
    if (btnAddLevel) btnAddLevel.addEventListener('click', addBlindLevel);

    var btnAddBreak = $('btnAddBreak');
    if (btnAddBreak) btnAddBreak.addEventListener('click', addBreakLevel);

    var btnCalcPrize = $('btnCalcPrize');
    if (btnCalcPrize) btnCalcPrize.addEventListener('click', calculatePrizes);

    var btnAddPrize = $('btnAddPrize');
    if (btnAddPrize) btnAddPrize.addEventListener('click', addPrizeRow);

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
