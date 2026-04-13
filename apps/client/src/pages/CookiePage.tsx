import { useState, useEffect } from "react";
import { getCookies, setCookies as apiSetCookies } from "../api/client";

export default function CookiePage() {
  const [cookieInput, setCookieInput] = useState("");
  const [configured, setConfigured] = useState<boolean | null>(null);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState<{ type: "success" | "error"; text: string } | null>(null);

  useEffect(() => {
    checkCookie();
  }, []);

  async function checkCookie() {
    try {
      const res = await getCookies();
      setConfigured(res.configured);
    } catch {
      setConfigured(false);
    }
  }

  async function handleSave() {
    setLoading(true);
    setMessage(null);
    try {
      await apiSetCookies(cookieInput);
      setConfigured(true);
      setMessage({ type: "success", text: "Cookie 保存成功" });
    } catch (e: any) {
      setMessage({ type: "error", text: e.message });
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      <div className="page-title">Cookie 配置</div>

      <div className="form-group">
        <label>小红书 Cookie</label>
        <textarea
          placeholder="粘贴你的小红书 Cookie..."
          value={cookieInput}
          onChange={(e) => setCookieInput(e.target.value)}
        />
      </div>

      <div className="form-group">
        <label>状态</label>
        <div style={{ padding: "8px 0", color: configured ? "#43a047" : "#e53935" }}>
          {configured === null ? "检查中..." : configured ? "已配置" : "未配置"}
        </div>
      </div>

      <div className="btn-group">
        <button className="btn btn-primary" onClick={handleSave} disabled={loading}>
          {loading ? "保存中..." : "保存"}
        </button>
      </div>

      {message && (
        <div className={message.type === "success" ? "success" : "error"} style={{ marginTop: 12 }}>
          {message.text}
        </div>
      )}
    </div>
  );
}
