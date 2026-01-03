---
id: 1911
title: "ポーカー"
slug: "poker"
status: publish
parent: 0
menu_order: 0
---

<style>
/* ポーカータイマー - Phase 6 スタイル */

/* 認証オーバーレイ - !important で上書き防止 */
.auth-overlay {
  position: fixed !important;
  top: 0 !important;
  left: 0 !important;
  width: 100% !important;
  height: 100% !important;
  background: rgba(0,0,0,0.9) !important;
  display: flex !important;
  justify-content: center !important;
  align-items: center !important;
  z-index: 2147483647 !important;
  margin: 0 !important;
  padding: 0 !important;
}

.auth-overlay.hidden {
  display: none !important;
}

.auth-modal {
  background: #fff !important;
  padding: 40px !important;
  border-radius: 16px !important;
  text-align: center !important;
  max-width: 400px !important;
  width: 90% !important;
  box-shadow: 0 20px 60px rgba(0,0,0,0.5) !important;
  position: relative !important;
  z-index: 2147483647 !important;
}

.auth-modal h2 {
  margin: 0 0 24px 0 !important;
  color: #1e293b !important;
  font-size: 20px !important;
  font-weight: 600 !important;
}

.auth-modal input[type="password"] {
  width: 100% !important;
  padding: 14px !important;
  font-size: 16px !important;
  border: 2px solid #e2e8f0 !important;
  border-radius: 8px !important;
  margin-bottom: 16px !important;
  box-sizing: border-box !important;
  background: #fff !important;
  color: #1e293b !important;
}

.auth-modal input[type="password"]:focus {
  outline: none !important;
  border-color: #3b82f6 !important;
}

.auth-modal button {
  width: 100% !important;
  padding: 14px !important;
  background: #3b82f6 !important;
  color: #fff !important;
  border: none !important;
  border-radius: 8px !important;
  font-size: 16px !important;
  font-weight: 600 !important;
  cursor: pointer !important;
  transition: background 0.2s !important;
}

.auth-modal button:hover {
  background: #2563eb !important;
}

.auth-error {
  color: #ef4444 !important;
  margin-top: 12px !important;
  font-size: 14px !important;
  min-height: 20px !important;
}

/* メインアプリ */
.poker-timer-app {
  max-width: 1100px;
  margin: 0 auto;
  font-family: 'Helvetica Neue', Arial, sans-serif;
  background: #f8fafc;
  color: #1e293b;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.1);
  position: relative;
}

.poker-timer-app * {
  box-sizing: border-box;
}

/* 3カラムレイアウト */
.timer-grid {
  display: grid;
  grid-template-columns: 200px 1fr 200px;
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
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 10px;
  margin-bottom: 10px;
}

.prize-total {
  font-size: 14px;
  color: #64748b;
  margin-bottom: 4px;
}

.prize-total-value {
  font-size: 22px;
  font-weight: bold;
  color: #d97706;
}

.prize-inmoney {
  font-size: 12px;
  color: #64748b;
  margin-top: 4px;
}

.prize-list {
  flex: 1;
  overflow: hidden;
  position: relative;
}

.prize-list-inner {
  position: absolute;
  width: 100%;
}

.prize-item {
  display: flex;
  justify-content: space-between;
  padding: 6px 0;
  font-size: 13px;
  border-bottom: 1px solid #f1f5f9;
}

.prize-item:last-child {
  border-bottom: none;
}

.prize-rank {
  color: #64748b;
  font-weight: 600;
}

.prize-amount {
  color: #1e293b;
  font-weight: bold;
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
}

/* レベル表示 */
.level-badge {
  display: inline-block;
  background: #3b82f6;
  color: #fff;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 10px;
}

.level-badge.break {
  background: #22c55e;
}

/* タイマー表示 */
.timer-time {
  font-size: 72px;
  font-weight: bold;
  font-family: 'Courier New', monospace;
  letter-spacing: 4px;
  margin: 15px 0;
  text-shadow: 0 2px 10px rgba(0,0,0,0.3);
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
  height: 6px;
  background: rgba(255,255,255,0.2);
  border-radius: 3px;
  margin: 15px 0;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #fbbf24, #f59e0b);
  transition: width 1s linear;
  border-radius: 3px;
}

/* ブラインド情報 - 視認性向上 */
.blind-info {
  margin: 10px 0;
}

.blind-current {
  font-size: 34px;
  font-weight: bold;
  color: #fbbf24;
  margin: 12px 0;
  text-shadow: 0 2px 8px rgba(251, 191, 36, 0.3);
}

.blind-ante {
  font-size: 18px;
  color: #94a3b8;
}

.blind-next {
  font-size: 14px;
  color: #94a3b8;
  margin-top: 15px;
  padding-top: 15px;
  border-top: 2px solid rgba(255,255,255,0.15);
}

.blind-next-label {
  font-size: 12px;
  color: #64748b;
  margin-bottom: 4px;
}

.blind-next-value {
  font-size: 22px;
  font-weight: 600;
  color: #e2e8f0;
  display: block;
  margin-top: 4px;
}

.next-ante {
  font-size: 14px;
  color: #94a3b8;
  margin-left: 8px;
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
  padding-bottom: 160px;
}

/* フルスクリーンボタン（右上固定） */
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

.info-card {
  background: #fff;
  border-radius: 12px;
  padding: 15px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

/* ラベルスタイル */
.panel-label {
  font-size: 11px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 6px;
}

.panel-value {
  font-size: 20px;
  font-weight: bold;
  color: #1e293b;
}

.panel-value.gold {
  color: #d97706;
}

.panel-value.large {
  font-size: 28px;
}

/* NEXT BREAK IN */
.break-card {
  background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
  color: #fff;
}

.break-card .panel-label {
  color: rgba(255,255,255,0.8);
}

.break-card .panel-value {
  color: #fff;
}

.break-card.hidden {
  display: none;
}

/* STACK カード */
.stack-card .stack-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 0;
}

.stack-card .stack-label {
  font-size: 12px;
  color: #64748b;
}

.stack-card .stack-value {
  font-size: 16px;
  font-weight: bold;
  color: #1e293b;
}

/* PLAYERS カード - 右下固定 */
.players-card {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
}

.players-card .players-display {
  font-size: 32px;
  font-weight: bold;
  color: #1e293b;
  margin: 8px 0;
}

.players-card .players-label {
  font-size: 11px;
  color: #94a3b8;
  margin-bottom: 10px;
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
  z-index: 2147483646;
  overflow-y: auto;
}

.settings-modal.active {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 20px;
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
  font-size: 15px;
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
  font-size: 14px;
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
  gap: 6px;
  padding: 8px 10px;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
  cursor: grab;
  transition: background 0.2s, transform 0.2s;
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
  width: 55px;
  padding: 6px;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  background: #fff;
  color: #1e293b;
  font-size: 13px;
  text-align: center;
}

.blind-level-item input.time-input {
  width: 45px;
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
  width: 24px;
  font-weight: bold;
  color: #3b82f6;
  font-size: 13px;
}

.blind-level-item .break-label {
  color: #22c55e;
  font-weight: 600;
  font-size: 13px;
}

.blind-level-item .delete-level {
  background: #ef4444;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 4px 8px;
  cursor: pointer;
  font-size: 12px;
  margin-left: auto;
}

.blind-level-item .time-label {
  font-size: 11px;
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
  gap: 6px;
  padding: 10px 12px;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
  cursor: grab;
  transition: background 0.2s;
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
  padding: 6px;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  font-size: 14px;
  text-align: center;
}

.prize-edit-item input.rank-input {
  width: 50px;
}

.prize-edit-item input.amount-input {
  width: 80px;
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
  font-size: 12px;
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
}

.poker-timer-app:fullscreen .timer-grid {
  flex: 1;
}

.poker-timer-app:fullscreen .timer-time {
  font-size: 120px;
}

.poker-timer-app:fullscreen .blind-current {
  font-size: 42px;
}

.poker-timer-app:fullscreen .blind-next-value {
  font-size: 28px;
}

/* レスポンシブ */
@media (max-width: 900px) {
  .timer-grid {
    grid-template-columns: 1fr;
  }

  .left-panel {
    order: 3;
    max-height: 200px;
  }

  .center-panel {
    order: 1;
  }

  .right-panel {
    order: 2;
    flex-direction: row;
    flex-wrap: wrap;
    padding-top: 0;
    padding-bottom: 0;
  }

  .right-panel .info-card {
    flex: 1;
    min-width: 140px;
  }

  .fullscreen-btn-top {
    position: static;
    margin-left: auto;
  }

  .players-card {
    position: static;
  }

  .timer-time {
    font-size: 56px;
  }

  .blind-current {
    font-size: 28px;
  }

  .prize-list {
    max-height: 120px;
  }
}

@media (max-width: 500px) {
  .poker-timer-app {
    padding: 12px;
  }

  .timer-time {
    font-size: 42px;
  }

  .blind-current {
    font-size: 22px;
  }

  .controls {
    flex-wrap: wrap;
  }

  .btn {
    padding: 10px 16px;
    font-size: 13px;
  }
}
</style>

<!-- 認証画面 - インラインスタイルで確実に表示 -->
<div class="auth-overlay" id="authOverlay" style="position:fixed!important;top:0!important;left:0!important;right:0!important;bottom:0!important;width:100vw!important;height:100vh!important;background:rgba(0,0,0,0.95)!important;display:flex!important;justify-content:center!important;align-items:center!important;z-index:2147483647!important;margin:0!important;padding:0!important;overflow:hidden!important;">
  <div class="auth-modal" style="background:#fff!important;padding:40px!important;border-radius:16px!important;text-align:center!important;max-width:400px!important;width:90%!important;box-shadow:0 20px 60px rgba(0,0,0,0.5)!important;position:relative!important;z-index:2147483647!important;margin:0 auto!important;">
    <h2 style="margin:0 0 24px 0!important;color:#1e293b!important;font-size:20px!important;font-weight:600!important;">パスワードを入力してください</h2>
    <input type="password" id="authPassword" placeholder="パスワード" style="width:100%!important;padding:14px!important;font-size:16px!important;border:2px solid #e2e8f0!important;border-radius:8px!important;margin-bottom:16px!important;box-sizing:border-box!important;background:#fff!important;color:#1e293b!important;">
    <button id="btnAuth" style="width:100%!important;padding:14px!important;background:#3b82f6!important;color:#fff!important;border:none!important;border-radius:8px!important;font-size:16px!important;font-weight:600!important;cursor:pointer!important;">認証</button>
    <div class="auth-error" id="authError" style="color:#ef4444!important;margin-top:12px!important;font-size:14px!important;min-height:20px!important;"></div>
  </div>
</div>

<!-- メインアプリ（初期非表示） -->
<div class="poker-timer-app" id="pokerTimer" style="display: none;">
  <div class="timer-grid">
    <!-- 左カラム: PRIZE -->
    <div class="left-panel">
      <div class="prize-header">
        <div class="panel-label">PRIZE</div>
        <div class="prize-total-value" id="prizeTotal">2,400 pt</div>
        <div class="prize-inmoney" id="prizeInmoney">インマネ: 2名</div>
      </div>
      <div class="prize-list" id="prizeListContainer">
        <div class="prize-list-inner" id="prizeList">
        </div>
      </div>
    </div>

    <!-- 中央カラム: タイマー -->
    <div class="center-panel">
      <div class="level-badge" id="levelBadge">Level <span id="currentLevel">1</span> / <span id="totalLevels">10</span></div>

      <div class="timer-time" id="timerDisplay">10:00</div>

      <div class="progress-bar">
        <div class="progress-fill" id="progressBar" style="width: 100%"></div>
      </div>

      <div class="blind-info">
        <div class="blind-current" id="currentBlind">25 / 50</div>
        <div class="blind-ante" id="anteInfo">Ante: 0</div>
        <div class="blind-next">
          <div class="blind-next-label">Next</div>
          <span class="blind-next-value" id="nextBlind">50 / 100 <span class="next-ante">(Ante: 0)</span></span>
        </div>
      </div>

      <div class="controls">
        <button class="btn btn-primary" id="startBtn">開始</button>
        <button class="btn btn-secondary" id="btnPrev">戻る</button>
        <button class="btn btn-secondary" id="btnSkip">スキップ</button>
        <button class="btn btn-warning" id="btnSettings">設定</button>
      </div>
    </div>

    <!-- 右カラム -->
    <div class="right-panel">
      <!-- フルスクリーンボタン（右上） -->
      <button class="fullscreen-btn-top" id="btnFullscreen" title="フルスクリーン">&#9974;</button>

      <!-- NEXT BREAK IN -->
      <div class="info-card break-card" id="breakCard">
        <div class="panel-label">NEXT BREAK IN</div>
        <div class="panel-value large" id="nextBreakTime">--:--</div>
      </div>

      <!-- STACK（中央） -->
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
          <button class="counter-btn" id="btnEntryMinus" title="エントリー減">&#9664;</button>
          <button class="counter-btn" id="btnRemainPlus" title="残り+">&#9650;</button>
          <button class="counter-btn" id="btnRemainMinus" title="残り-">&#9660;</button>
          <button class="counter-btn" id="btnEntryPlus" title="エントリー増">&#9654;</button>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- 設定モーダル -->
<div class="settings-modal" id="settingsModal">
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
          <label class="setting-label">デフォルトの時間（分）</label>
          <input type="number" class="setting-input" id="defaultLevelTime" value="10" min="1" max="60">
          <small style="color: #64748b; font-size: 11px;">新規レベル追加時のデフォルト値</small>
        </div>
        <div class="setting-group">
          <label class="setting-label">初期エントリー数</label>
          <input type="number" class="setting-input" id="initialPlayers" value="8" min="2" max="100">
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
          <label class="setting-label">エントリーフィー (pt)</label>
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
  // パスワードハッシュ（shiimandy）
  var PASSWORD_HASH = 'fc35ac7c2ccf7cd4a6ca8f9fe15955385963a411acbb498a2e2da7221daad9b4';

  // デフォルトブラインド構成（レベルごとの時間付き）
  var defaultBlinds = [
    { sb: 25, bb: 50, ante: 0, isBreak: false, levelTime: 10 },
    { sb: 50, bb: 100, ante: 0, isBreak: false, levelTime: 10 },
    { sb: 75, bb: 150, ante: 0, isBreak: false, levelTime: 10 },
    { sb: 100, bb: 200, ante: 25, isBreak: false, levelTime: 10 },
    { sb: 150, bb: 300, ante: 25, isBreak: false, levelTime: 10 },
    { isBreak: true, breakTime: 5 },
    { sb: 200, bb: 400, ante: 50, isBreak: false, levelTime: 10 },
    { sb: 300, bb: 600, ante: 75, isBreak: false, levelTime: 10 },
    { sb: 400, bb: 800, ante: 100, isBreak: false, levelTime: 10 },
    { sb: 500, bb: 1000, ante: 100, isBreak: false, levelTime: 10 }
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
    defaultLevelTime: 10,
    totalPlayers: 8,
    remainingPlayers: 8,
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

  // SHA-256ハッシュ計算
  function sha256(text) {
    return new Promise(function(resolve, reject) {
      var encoder = new TextEncoder();
      var data = encoder.encode(text);
      crypto.subtle.digest('SHA-256', data).then(function(hashBuffer) {
        var hashArray = Array.from(new Uint8Array(hashBuffer));
        var hashHex = hashArray.map(function(b) { return b.toString(16).padStart(2, '0'); }).join('');
        resolve(hashHex);
      }).catch(reject);
    });
  }

  // 認証チェック
  function checkAuth() {
    if (sessionStorage.getItem('pokerTimerAuth') === 'true') {
      showApp();
      return;
    }
    showAuthScreen();
  }

  // 認証実行
  function authenticate() {
    var input = $('authPassword').value;
    sha256(input).then(function(hash) {
      if (hash === PASSWORD_HASH) {
        sessionStorage.setItem('pokerTimerAuth', 'true');
        showApp();
      } else {
        $('authError').textContent = 'パスワードが違います';
      }
    });
  }

  // アプリ表示
  function showApp() {
    var overlay = $('authOverlay');
    if (overlay) {
      overlay.style.cssText = 'display:none!important;';
    }
    var app = $('pokerTimer');
    if (app) app.style.display = 'block';
    loadSettings();
  }

  // 認証画面表示
  function showAuthScreen() {
    var overlay = $('authOverlay');
    if (overlay) {
      overlay.style.cssText = 'position:fixed!important;top:0!important;left:0!important;right:0!important;bottom:0!important;width:100vw!important;height:100vh!important;background:rgba(0,0,0,0.95)!important;display:flex!important;justify-content:center!important;align-items:center!important;z-index:2147483647!important;margin:0!important;padding:0!important;overflow:hidden!important;';
    }
    var app = $('pokerTimer');
    if (app) app.style.display = 'none';
  }

  // LocalStorageから設定を読み込み
  function loadSettings() {
    var saved = localStorage.getItem('pokerTimerSettingsV7');
    if (saved) {
      try {
        var settings = JSON.parse(saved);
        state.defaultLevelTime = settings.defaultLevelTime || 10;
        state.totalPlayers = settings.totalPlayers || 8;
        state.remainingPlayers = settings.totalPlayers || 8;
        state.initialStack = settings.initialStack || 500;
        state.blinds = settings.blinds || JSON.parse(JSON.stringify(defaultBlinds));
        state.blindSets = settings.blindSets || { 'default': JSON.parse(JSON.stringify(defaultBlinds)) };
        state.currentSetName = settings.currentSetName || 'default';
        state.prizeSettings = settings.prizeSettings || JSON.parse(JSON.stringify(defaultPrizeSettings));
        state.prizeDistribution = settings.prizeDistribution || [];

        // 現在レベルの時間を設定
        var currentBlind = state.blinds[state.currentLevel];
        if (currentBlind) {
          if (currentBlind.isBreak) {
            state.levelTime = (currentBlind.breakTime || 5) * 60;
          } else {
            state.levelTime = (currentBlind.levelTime || state.defaultLevelTime) * 60;
          }
          state.timeRemaining = state.levelTime;
        }
      } catch(e) {
        console.log('Settings load error:', e);
      }
    }
    updateBlindSetSelector();
    updateDisplay();
    updatePrizeDisplay();
    startPrizeAutoScroll();
  }

  // 設定を保存
  function saveSettings() {
    var defaultLevelTime = parseInt($('defaultLevelTime').value) || 10;
    var totalPlayers = parseInt($('initialPlayers').value) || 8;
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

    localStorage.setItem('pokerTimerSettingsV7', JSON.stringify({
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
    var defaultTime = parseInt($('defaultLevelTime').value) || 10;
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
    state.totalPlayers = Math.max(1, state.totalPlayers + delta);
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
    var cb = $('currentBlind');
    var ai = $('anteInfo');
    var nb = $('nextBlind');

    if (blind) {
      if (blind.isBreak) {
        if (lb) {
          lb.textContent = 'BREAK';
          lb.classList.add('break');
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

    var showBreak = false;
    if (foundBreak) {
      if (!state.blinds[state.currentLevel].isBreak) {
        showBreak = true;
      }
    }
    if (showBreak) {
      breakCard.classList.remove('hidden');
      var minutes = Math.floor(timeToBreak / 60);
      var seconds = timeToBreak % 60;
      var nbt = $('nextBreakTime');
      if (nbt) nbt.textContent = String(minutes).padStart(2, '0') + ':' + String(seconds).padStart(2, '0');
    } else {
      breakCard.classList.add('hidden');
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

    var pt = $('prizeTotal');
    if (pt) pt.textContent = Math.round(P).toLocaleString() + ' pt';
    var pi = $('prizeInmoney');
    if (pi) pi.textContent = 'インマネ: ' + totalInMoney + '名';

    var listEl = $('prizeList');
    if (listEl) {
      listEl.innerHTML = '';

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

    if (inner.offsetHeight > container.offsetHeight) {
      inner.style.animation = 'none';
      startPrizeAutoScroll();
    } else {
      if (state.scrollIntervalId) {
        clearInterval(state.scrollIntervalId);
        state.scrollIntervalId = null;
      }
      inner.style.top = '0';
    }
  }

  // 片方向ループスクロール
  function startPrizeAutoScroll() {
    if (state.scrollIntervalId) {
      clearInterval(state.scrollIntervalId);
    }

    var container = $('prizeListContainer');
    var inner = $('prizeList');

    if (!inner || !container || inner.offsetHeight <= container.offsetHeight) return;

    var scrollPos = 0;
    var totalHeight = inner.offsetHeight;
    var containerHeight = container.offsetHeight;
    var pauseCounter = 0;

    state.scrollIntervalId = setInterval(function() {
      if (pauseCounter > 0) {
        pauseCounter--;
        return;
      }

      scrollPos += 1;

      if (scrollPos >= totalHeight) {
        scrollPos = -containerHeight;
        pauseCounter = 30;
      }

      inner.style.top = -scrollPos + 'px';
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

    if (!document.fullscreenElement) {
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
    // 認証
    var btnAuth = $('btnAuth');
    if (btnAuth) btnAuth.addEventListener('click', authenticate);

    var authPassword = $('authPassword');
    if (authPassword) {
      authPassword.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') authenticate();
      });
    }

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
    checkAuth();
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
</script>
