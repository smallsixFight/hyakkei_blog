const REQUEST_URL = "http://127.0.0.1:9900";

// 路由
const api = {
    articles: `/api/v1/articles`,
    books: `/api/v1/books`,
    friends: `/api/v1/friends`,
    visitorAdd: `/api/v1/visitor/add`,
}

// 消息提示
let message = new Object({
    count: 0,
    timeoutIds: [],
});

message.destory = function () {
    $('.message').remove();
    this.count = 0;
    for (let i = 0; i < this.timeoutIds; i++) {
        this.removeTimeoutEvent(this.timeoutIds[i]);
    }
}

message.showInfo = function (msg, duration, typ) {
    typ = typ === 'success' ? 'success' : typ === 'warn' ? 'warn' : typ === 'error' ? 'error' : '';
    typ = typ ? 'message-' + typ : '';
    msg = msg || '';
    duration = (isNaN(duration) || duration < 500) ? 3000 : duration;
    const id = "msg-" + new Date().getTime();
    if (this.count) {
        let list = $('.message');
        for (let i = 0; i < list.length; i++) {
            const topDist = list[i].style.top;
            list[i].style.top = (parseInt(topDist) + 64) + 'px';
        }
    }
    let node = '<div id="' + id + '" class="message ' + typ + ' message-enter" style="top: 20px;"><p>' + msg || '' + '</p></div>';
    this.count++;
    $('body').prepend(node);
    setTimeout(() => $('#' + id).removeClass('message-enter'), 10);
    const tid = setTimeout(() => {
        this.removeMessage(id);
        this.removeTimeoutEvent(tid)
    }, parseInt(duration));
    this.timeoutIds.push(tid);
}

message.removeMessage = function (id) {
    $('#' + id).addClass('message-leave-out');
    const tid = setTimeout(() => {
        $('#' + id).remove();
        this.count -= 1;
        this.removeTimeoutEvent(tid);
    }, 400);
    this.timeoutIds.push(tid);
}

message.removeTimeoutEvent = function (id) {
    const idx = this.timeoutIds.indexOf(id);
    if (idx > -1) {
        this.timeoutIds.splice(idx, 1);
    }
}

message.info = function (msg, duration) {
    this.showInfo(msg, duration);
}

message.success = function (msg, duration) {
    this.showInfo(msg, duration, 'success');
}

message.warn = function (msg, duration) {
    this.showInfo(msg, duration, 'warn');
}

message.error = function (msg, duration) {
    this.showInfo(msg, duration, 'error');
}

async function uploadFiles(path, fileData, params) {
    let result;
    let uploadData = new FormData();
    uploadData.append("file", fileData);
    for (let i in params) {
        uploadData.append(i, params[i]);
    }
    await $.ajax({
        url: REQUEST_URL + path,
        type: 'POST',
        cache: false,
        data: uploadData,
        timeout: 90000,
        processData: false,
        contentType: false,
        dataType: "json",
        success: function (resp) {
            result = resp;
        },
        beforeSend: (xhr) => beforeSend(xhr),
        error: (xhr) => handleError(xhr),
    });
    return result;
}

async function request(path, params = {data: null, method: "GET", contentType: "application/json;charset=utf-8"}) {
    let result;
    await $.ajax({
        url: REQUEST_URL + path,
        data: params.data && JSON.stringify(params.data),
        type: params.method,
        timeout: 30000,
        contentType: params.contentType,
        beforeSend: (xhr) => beforeSend(xhr),
        success: function (resp) {
            result = resp;
        },
        error: (xhr) => handleError(xhr),
    })
    return result;
}

function beforeSend(xhr) {
    const token = sessionStorage.getItem("token");
    if (token) {
        xhr.setRequestHeader("accessToken", token);
    }
}

function handleError(xhr) {
    let msg = xhr.responseJSON && (xhr.responseJSON.message || xhr.status);
    message.error("请求失败：" + msg);
}

function getqueryParam() {
    let str = window.location.search;
    if (str.length < 1) {
        return;
    }
    str = str.substr(1)
    const arr = str.split("&")
    const params = {};
    for (let i in arr) {
        const kv = arr[i].split("=");
        kv.length === 2 && (params[kv[0]] = kv[1]);
    }
    return params;
}

request(api.visitorAdd);