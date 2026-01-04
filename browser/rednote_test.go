package browser

import (
	"context"
	"testing"
)

func TestExtractRednoteUrlAndTitle(t *testing.T) {
	ctx := context.Background()

	// 测试用例1：正常HTML内容（使用真实小红书格式）
	testHTML := `
<!DOCTYPE html>
<html>
<body>
	<section class="note-item">
		<div class="footer">
			<a class="title"><span>12306候补20%能候补成功吗！</span></a>
		</div>
		<a href="/explore/695a0fb9000000001d03a898" style="display: none;"></a>
	</section>
	<section class="note-item">
		<div class="footer">
			<a class="title"><span>测试标题2</span></a>
		</div>
		<a href="/explore/note2" style="display: none;"></a>
	</section>
</body>
</html>`

	result, err := ExtractRednoteUrlAndTitle(ctx, testHTML)
	if err != nil {
		t.Fatalf("正常测试用例失败：%v", err)
	}

	if len(result) != 2 {
		t.Errorf("期望提取到2个笔记，实际得到：%d", len(result))
	}

	// 验证第一个笔记
	if result[0].Title != "12306候补20%能候补成功吗！" {
		t.Errorf("第一个笔记标题期望：12306候补20%%能候补成功吗！，实际：%s", result[0].Title)
	}
	if result[0].Url != "https://www.xiaohongshu.com/explore/695a0fb9000000001d03a898" {
		t.Errorf("第一个笔记URL期望：https://www.xiaohongshu.com/explore/695a0fb9000000001d03a898，实际：%s", result[0].Url)
	}

	// 验证第二个笔记
	if result[1].Title != "测试标题2" {
		t.Errorf("第二个笔记标题期望：测试标题2，实际：%s", result[1].Title)
	}
	if result[1].Url != "https://www.xiaohongshu.com/explore/note2" {
		t.Errorf("第二个笔记URL期望：https://www.xiaohongshu.com/explore/note2，实际：%s", result[1].Url)
	}

	// 测试用例2：空HTML内容
	_, err = ExtractRednoteUrlAndTitle(ctx, "")
	if err == nil {
		t.Error("空HTML内容应该返回错误")
	}

	// 测试用例3：只有空白字符的HTML内容
	_, err = ExtractRednoteUrlAndTitle(ctx, "   \n\t  ")
	if err == nil {
		t.Error("只有空白字符的HTML内容应该返回错误")
	}

	// 测试用例4：无效的HTML
	invalidHTML := "<html><body><invalid>content</invalid>"
	_, err = ExtractRednoteUrlAndTitle(ctx, invalidHTML)
	if err == nil {
		t.Error("无效的HTML应该返回错误")
	}

	// 测试用例5：没有有效笔记项的HTML
	noNotesHTML := `
<!DOCTYPE html>
<html>
<body>
	<section class="other-class">
		<a href="/explore/note1">标题</a>
	</section>
	<section class="query-note-wrapper">
		<a href="/explore/note2">标题</a>
	</section>
</body>
</html>`

	_, err = ExtractRednoteUrlAndTitle(ctx, noNotesHTML)
	if err == nil {
		t.Error("没有有效笔记项的HTML应该返回错误")
	}

	// 测试用例6：笔记项没有URL的HTML
	noURLHTML := `
<!DOCTYPE html>
<html>
<body>
	<section class="note-item">
		<div class="footer">
			<a class="title"><span>测试标题</span></a>
		</div>
	</section>
</body>
</html>`

	result, err = ExtractRednoteUrlAndTitle(ctx, noURLHTML)
	if err != nil {
		t.Fatalf("没有URL的HTML测试用例失败：%v", err)
	}
	if len(result) != 0 {
		t.Errorf("没有有效URL的笔记项应该被跳过，期望：0，实际：%d", len(result))
	}

	// 测试用例7：笔记项没有标题的HTML
	noTitleHTML := `
<!DOCTYPE html>
<html>
<body>
	<section class="note-item">
		<a href="/explore/note1" style="display: none;"></a>
	</section>
</body>
</html>`

	result, err = ExtractRednoteUrlAndTitle(ctx, noTitleHTML)
	if err != nil {
		t.Fatalf("没有标题的HTML测试用例失败：%v", err)
	}
	if len(result) != 1 {
		t.Errorf("期望提取到1个笔记，实际得到：%d", len(result))
	}
	if result[0].Title != "" {
		t.Errorf("没有标题的笔记应该返回空字符串，实际：%s", result[0].Title)
	}
	if result[0].Url != "https://www.xiaohongshu.com/explore/note1" {
		t.Errorf("URL应该正确生成，实际：%s", result[0].Url)
	}

	// 测试用例8：包含特殊字符的标题
	specialCharHTML := `
<!DOCTYPE html>
<html>
<body>
	<section class="note-item">
		<div class="footer">
			<a class="title"><span>测试 &amp; 标题 &lt;特殊&gt; 字符</span></a>
		</div>
		<a href="/explore/note1" style="display: none;"></a>
	</section>
</body>
</html>`

	result, err = ExtractRednoteUrlAndTitle(ctx, specialCharHTML)
	if err != nil {
		t.Fatalf("特殊字符标题测试用例失败：%v", err)
	}
	if len(result) != 1 {
		t.Errorf("期望提取到1个笔记，实际得到：%d", len(result))
	}
	// HTML实体应该被正确解码
	if result[0].Title != "测试 & 标题 <特殊> 字符" {
		t.Errorf("特殊字符标题解码错误，期望：测试 & 标题 <特殊> 字符，实际：%s", result[0].Title)
	}

	// 测试用例9：使用真实小红书HTML结构
	realXiaohongshuHTML := `
<!DOCTYPE html>
<html>
<body>
	<div class="search-layout__main">
		<div class="feeds-container">
			<section class="note-item" data-width="1200" data-height="1600" data-index="0">
				<a href="/explore/695a0fb9000000001d03a898" style="display: none;"></a>
				<div class="footer">
					<a class="title">
						<span>12306候补20%能候补成功吗！</span>
					</a>
					<div class="card-bottom-wrapper">
						<a class="author">
							<div class="name-time-wrapper">
								<div class="name">莉莉日记</div>
								<div class="time">1小时前</div>
							</div>
						</a>
						<span class="like-wrapper">
							<span class="count">赞</span>
						</span>
					</div>
				</div>
			</section>
		</div>
	</div>
</body>
</html>`

	result, err = ExtractRednoteUrlAndTitle(ctx, realXiaohongshuHTML)
	if err != nil {
		t.Fatalf("真实小红书HTML测试用例失败：%v", err)
	}
	if len(result) != 1 {
		t.Errorf("期望提取到1个笔记，实际得到：%d", len(result))
	}
	if result[0].Title != "12306候补20%能候补成功吗！" {
		t.Errorf("真实HTML标题提取错误，期望：12306候补20%%能候补成功吗！，实际：%s", result[0].Title)
	}
	if result[0].Url != "https://www.xiaohongshu.com/explore/695a0fb9000000001d03a898" {
		t.Errorf("真实HTML URL提取错误，期望：https://www.xiaohongshu.com/explore/695a0fb9000000001d03a898，实际：%s", result[0].Url)
	}
}

func TestExtractRednoteUrlAndTitle_EdgeCases(t *testing.T) {
	ctx := context.Background()

	// 测试用例：正常URL但标题有空格
	spaceTitleHTML := `
<!DOCTYPE html>
<html>
<body>
	<section class="note-item">
		<div class="footer">
			<a class="title"><span>  测试标题  </span></a>
		</div>
		<a href="/explore/note1" style="display: none;"></a>
	</section>
</body>
</html>`

	result, err := ExtractRednoteUrlAndTitle(ctx, spaceTitleHTML)
	if err != nil {
		t.Fatalf("标题空格测试用例失败：%v", err)
	}
	if len(result) != 1 {
		t.Errorf("期望提取到1个笔记，实际得到：%d", len(result))
	}
	// 标题中的空格应该被trim掉
	if result[0].Title != "测试标题" {
		t.Errorf("标题空格应该被trim，期望：测试标题，实际：%s", result[0].Title)
	}

	// 测试用例：混合有效和无效的笔记项
	mixedHTML := `
<!DOCTYPE html>
<html>
<body>
	<section class="note-item">
		<div class="footer">
			<a class="title"><span>有效标题1</span></a>
		</div>
		<a href="/explore/note1" style="display: none;"></a>
	</section>
	<section class="note-item">
		<div class="footer">
			<a class="title"><span>无效标题</span></a>
		</div>
		<!-- 没有URL -->
	</section>
	<section class="note-item">
		<div class="footer">
			<a class="title"><span>有效标题2</span></a>
		</div>
		<a href="/explore/note2" style="display: none;"></a>
	</section>
</body>
</html>`

	result, err = ExtractRednoteUrlAndTitle(ctx, mixedHTML)
	if err != nil {
		t.Fatalf("混合内容测试用例失败：%v", err)
	}
	// 只有2个有效笔记项（有URL的）
	if len(result) != 2 {
		t.Errorf("期望提取到2个有效笔记，实际得到：%d", len(result))
	}
	if result[0].Title != "有效标题1" || result[1].Title != "有效标题2" {
		t.Errorf("标题提取错误，实际：%v", []string{result[0].Title, result[1].Title})
	}
}
