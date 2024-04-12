// 获取计数器元素和按钮元素
var counterElement = document.getElementById("counter");
var incrementButton = document.getElementById("incrementBtn");

// 初始化计数器值
var count = 0;

// 点击按钮时增加计数器的值
incrementButton.addEventListener("click", function() {
  count++; // 增加计数器的值
  counterElement.innerText = count; // 更新计数器元素的文本内容
});