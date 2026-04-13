import { useState } from "react";
import { collectSingleNote } from "../api/client";

export default function SingleNotePage() {
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);

  async function handleCollect() {
    if (!url.trim()) {
      setError("请输入笔记 URL");
      return;
    }
    setLoading(true);
    setError(null);
    setResult(null);
    try {
      const data = await collectSingleNote(url);
      setResult(data);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      <div className="page-title">单条笔记采集</div>

      <div className="form-group">
        <label>笔记 URL</label>
        <input
          type="text"
          placeholder="https://www.xiaohongshu.com/explore/..."
          value={url}
          onChange={(e) => setUrl(e.target.value)}
        />
      </div>

      <div className="btn-group">
        <button className="btn btn-primary" onClick={handleCollect} disabled={loading}>
          {loading ? "采集中..." : "采集"}
        </button>
      </div>

      {error && <div className="error" style={{ marginTop: 12 }}>{error}</div>}

      {result && (
        <div className="result-box">
          <div className="title">采集结果</div>
          <div className="info">
            <p><strong>标题：</strong>{result.title || result.display_title || "-"}</p>
            <p><strong>作者：</strong>{result.user?.nickname || result.author?.name || "-"}</p>
            <p><strong>点赞：</strong>{result.interaction?.stat?.liked_count || "-"}</p>
            <p><strong>收藏：</strong>{result.interaction?.stat?.collected_count || "-"}</p>
            <p><strong>评论：</strong>{result.interaction?.stat?.comment_count || "-"}</p>
          </div>
        </div>
      )}
    </div>
  );
}
