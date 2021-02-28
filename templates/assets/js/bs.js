
// 路由
const bsApi = {
    login: `/api/v1/login`,
    dashboard: `/api/v1/dashboard/info`,
    bookList: `/api/v1/book/list`,
    saveBook: `/api/v1/book/save`,
    delBook: `/api/v1/book/del`,
    friendList: `/api/v1/friend/list`,
    saveFriendInfo: `/api/v1/friend/save`,
    delFriendLink: `/api/v1/friend/del`,
    tagList: `/api/v1/tag/list`,
    saveTagInfo: `/api/v1/tag/save`,
    delTag: `/api/v1/tag/del`,
    articleList: `/api/v1/article/list`,
    articleDetail: `/api/v1/article/detail`,
    delArticle: `/api/v1/article/del`,
    saveArticle: `/api/v1/article/save`,
    pageList: `/api/v1/page/list`,
    pageDetail: `/api/v1/page/detail`,
    delPage: `/api/v1/page/del`,
    savePage: `/api/v1/page/save`,
    postPreview: `/api/v1/post/markdown/preview`,
    getSysSetting: `/api/v1/sys_setting/info`,
    updateSysSetting: `/api/v1/sys_setting/update`,
    uploadFile: `/api/v1/upload/file`,

}

$('#username').text(sessionStorage.getItem("username") || 'hyakkei');

// 上传组件监听
$(".file-upload button").bind("click", function () {
    $(this).parent().find('input').click();
})

// 菜单弹出/隐藏
$('#collapsed-btn, .cover').click(() => {
    if ($('.backstage-aside').hasClass('show')) {
        $('.backstage-aside').removeClass('show');
        $('.cover').get(0).style.visibility = 'hidden';
    } else {
        $('.backstage-aside').addClass('show');
        $('.cover').get(0).style.visibility = 'visible';
    }
});

// 页面跳转
function jump(pathname, queryParams) {
    let p = "?";
    for (k in queryParams) {
        p += `${k}=${queryParams[k]}&`
    }
    window.location.href = pathname + '.html' + p.substr(0, p.length - 1);
}

// 文本编辑器监听
$('.edit-navigation').find('div').click((e) => {
    const sp = $('#articleContent').prop("selectionStart");
    const ep = $('#articleContent').prop("selectionEnd");
    let preContent = $('#articleContent').val().substring(0, sp);
    let afterContent = $('#articleContent').val().substring(ep);
    const content = $('#articleContent').val().substring(sp, ep) || '';
    let newPointer = 0;
    switch ($(e.currentTarget).attr('id')) {
        case "bold":
            $("#articleContent").val(preContent + "**" + $('#articleContent').val().substring(sp, ep) + "**" + afterContent);
            newPointer = ep + 2;
            break;
        case "italic":
            $("#articleContent").val(preContent + "*" + $('#articleContent').val().substring(sp, ep) + "*" + afterContent);
            newPointer = ep + 1;
            break;
        case "underline":
            $("#articleContent").val(preContent + "++" + $('#articleContent').val().substring(sp, ep) + "++" + afterContent);
            newPointer = ep + 2;
            break;
        case "strikethrough":
            $("#articleContent").val(preContent + "~~" + $('#articleContent').val().substring(sp, ep) + "~~" + afterContent);
            newPointer = ep + 2;
            break;
        case "quote":
            $("#articleContent").val(preContent + "> " + $('#articleContent').val().substring(sp, ep) + afterContent);
            newPointer = ep + 2;
            break;
        case "hr":
            $("#articleContent").val(preContent + "\n* * *\n\n" + $('#articleContent').val().substring(sp, ep) + afterContent);
            newPointer = sp + 8;
            break;
        case "tag_code":
            $("#articleContent").val(preContent + "`" + $('#articleContent').val().substring(sp, ep) + '`' + afterContent);
            newPointer = ep + 1;
            break;
        case "code_block":
            $("#articleContent").val(preContent + "\n```\n" + $('#articleContent').val().substring(sp, ep) + '\n```\n\n' + afterContent);
            newPointer = sp === ep ? sp + 5 : (2 * ep - sp + 6);
            break;
        case "link":
            if (content.startsWith("http")) {
                $("#articleContent").val(preContent + "[](" + content + ")\n" + afterContent);
                newPointer = sp + 1;
            } else {
                $("#articleContent").val(preContent + "[" + content + "]()\n" + afterContent);
                newPointer = sp + 3 + content.length;
            }
            break;
        case "bulletedlist":
            var newContent = preContent ? '\n- ' : '- ';
            var moveDist = newContent.length;
            for (let i = 0; i < content.length; i++) {
                newContent += content[i];
                if (content[i] === '\n') {
                    newContent += "- ";
                    moveDist += 2;
                }
            }
            $("#articleContent").val(preContent + newContent + afterContent);
            newPointer = ep + moveDist;
            break;
        case "numberedlist":
            var newContent = preContent ? '\n1. ' : '1. ';
            var moveDist = newContent.length;
            let num = 1;
            for (let i = 0; i < content.length; i++) {
                newContent += content[i];
                if (content[i] === '\n') {
                    num++;
                    newContent += num + ". ";
                    moveDist += 3;
                }
            }
            $("#articleContent").val(preContent + newContent + afterContent);
            newPointer = ep + moveDist;
            break;
        case "image":
            if (content.startsWith("http")) {
                $("#articleContent").val(preContent + "![](" + content + ")\n" + afterContent);
                newPointer = sp + 2;
            } else {
                $("#articleContent").val(preContent + "![" + content + "]()\n" + afterContent);
                newPointer = sp + 4 + content.length;
            }
            break;
    }
    $('#articleContent').focus();
    $('#articleContent').prop("selectionStart", newPointer);
    $('#articleContent').prop("selectionEnd", newPointer);
});

//  菜单导航选项更新
function onSelected(id) {
    $("#" + id).addClass("is-selected");
}

function generateTableContent(ele, columns) {
    if (!ele) {
        return;
    }
    let content = "<thead><tr>";
    for (let i in columns) {
        content += "<th>" + (columns[i].title || '') + "</th>"
    }
    content += "</tr></thead><tbody></tbody>";
    $(ele).append(content);
}

function updateTableData(ele, columns, dataList) {
    let bodyEle = $(ele).find("tbody")[0];
    if (!bodyEle) {
        return;
    }
    let rows = "";
    for (let i in dataList) {
        let tr = "<tr>";
        for (let j = 0; j < columns.length; j++) {
            const v = getVal(dataList[i], columns[j].dataIdx);
            tr += "<td>" + (columns[j].render ? columns[j].render(v, i) : v) + "</td>";
        }
        rows += tr + "</tr>";
    }
    $(bodyEle).empty();
    $(bodyEle).append(rows);
}

function getVal(info, key) {
    for (let k in info) {
        if (k === key) {
            return info[k];
        }
    }
    return "";
}

// 加载状态监听
function setLoading() {
    const eles = $(document).find("[my-loading]");
    eles.addClass("is-loading");
    for (let i = 0; i < eles.length; i++) {
        if (eles[i].tagName.toLowerCase() === 'button') {
            $(eles[i]).prepend("<img class='btn-loading-icon' src='../assets/img/icons/loading_btn.svg' />");
        } else {
            $(eles[i]).prepend("<img class='loading-icon' src='../assets/img/icons/loading.svg' />");
        }
    }
}

function removeLoading() {
    const eles = $(document).find("[my-loading]");
    eles.removeClass("is-loading");
    $(eles).children('.btn-loading-icon').remove();
    $(eles).children('.loading-icon').remove();
}

function logout() {
    sessionStorage.clear();
    jump("login")
}

// 设置选择框值
function setSelectVal(e, val) {
    let children = e.children();
    let hasOption = false;
    e.children().hasClass('is-selected') && e.children().removeClass('is-selected')
    for (let i = 0; i < children.length; i++) {
        if ($(children[i]).attr('value') === val) {
            e = $(children[i]);
            hasOption = true;
            break;
        }
    }
    if (!hasOption) { return };
    e.addClass('is-selected');
    const inputEle = $(e.parent().parent().find('input'));
    if (inputEle) {
        inputEle.val(e.text());
        inputEle.attr('value', e.attr('value') || e.text());
    }
}

// 选择框监听
$('.select li').bind("click", function () {
    let e = $(this);
    let children = e.parent().children();
    let isSelected = false;
    for (let i = 0; i < children.length; i++) {
        if ($(children[i]).hasClass('is-selected')) {
            isSelected = true;
            $(children[i]).removeClass('is-selected');
            break;
        }
    }
    // 如果是初始化，查看是否有 defalutValue，有则使用 defalutValue 设置的选项，否则使用当前选项（即第一个 li）
    if (!isSelected && e.parent().attr("defaultValue")) {
        const dv = e.parent().attr("defaultValue");
        for (let i = 0; i < children.length; i++) {
            if ($(children[i]).attr('value') === dv) {
                e = $(children[i]);
                break;
            }
        }
    }
    e.addClass('is-selected');
    const inputEle = $(e.parent().parent().find('input'));
    if (inputEle) {
        inputEle.val(e.text());
        inputEle.attr('value', e.attr('value') || e.text());
    }
    e.parent().attr("onchange") && window[e.parent().attr("onchange")](e.attr('value') || e.text());
});

// 多选框
function setMultSelectVal(e, vals) {
    let children = e.children();
    e.children().hasClass('is-selected') && e.children().removeClass('is-selected');
    const tagsDiv = e.parent().find('.multiple-select-tags').get(0);
    for (let i = 0; i < children.length; i++) {
        if (vals.indexOf($(children[i]).attr('value')) > -1) {
            $(children[i]).addClass('is-selected');
            $(tagsDiv).append(`<span value="${$(children[i]).attr('value')}">${$(children[i]).text()}</span>`);
        }
    }
    updateMultSelectHeight(tagsDiv);
    e.attr("onchange") && window[e.attr("onchange")](vals);
}

function updateMultSelectHeight(e) {
    const tags = $(e).children();
    const inputEle = $(e).parent().find('input');
    let spanToalWidth = 0;
    for (let tag of tags) {
        spanToalWidth += $(tag).outerWidth();
    }
    const v = spanToalWidth == 0 ? '请选择' : '';
    $(inputEle).attr('placeholder', v);
    inputEle.css('height', (Math.floor((spanToalWidth / inputEle.outerWidth() / 0.9)) * 40 + 40) + 'px');
}

function listenMultSelectClick() {
    $('.multiple-select li').bind("click", function () {
        let e = $(this);
        let vals = [];
        const tagsDiv = e.parent().parent().find('.multiple-select-tags').get(0);
        if (e.hasClass('is-selected')) {
            e.removeClass('is-selected');
            $(tagsDiv).find(`span[value=${e.attr('value')}]`).remove();
        } else {
            e.addClass('is-selected');
            $(tagsDiv).append(`<span value="${e.attr('value')}">${e.text()}</span>`);
        }
        const tags = $(tagsDiv).children();
        for (let tag of tags) {
            vals.push($(tag).attr('value'));
        }
        updateMultSelectHeight(tagsDiv);
        e.parent().attr("onchange") && window[e.parent().attr("onchange")](vals);
    });
}
listenMultSelectClick();

// 初始化
(function init() {
    selectedMenu();
    this.setTimeout(() => {
        setSelectDefalutVal();
    }, 300);
})()

// 菜单选择
function selectedMenu() {
    let pathname = document.location.pathname;
    const menuName = pathname.substring(pathname.lastIndexOf('/') + 1).split(".")[0].split("_")[0];
    const e = $(".backstage-aside");
    e && $("#aside-" + menuName).addClass("is-selected");
}

// 展示弹窗
function changeLayerVisible() {
    if ($(".global-layer").hasClass("layer-show")) {
        $(".global-layer").removeClass("layer-show");
        changeGlobalCover('none');
    } else {
        changeGlobalCover('block');
        $(".global-layer").addClass("layer-show");
    }
}

// 全局遮罩层
function changeGlobalCover(val) {
    $('.global-cover').css('display', val);
}

// 设置选择框默认选项
function setSelectDefalutVal() {
    const uls = document.querySelectorAll('.select ul');
    for (let ul of uls) {
        if ($(ul).children().hasClass("is-selected")) { return };
        const li = ul.querySelector('li');
        li && li.click();
    }
}