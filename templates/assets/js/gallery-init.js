(function (gallerySelector, itemSelector) {
    // 打开图片预览 gallery
    let openPhotoSwipe = function (galleryIndex, itemIndex) {
        let gallery = this.gallerys[galleryIndex];
        let options = {
            galleryUID: String(galleryIndex + 1),
            index: itemIndex,
            getThumbBoundsFn: function (index) {
                let thumbnail = gallery[index].el,
                    pageYScroll = window.pageYOffset || document.documentElement.scrollTop,
                    rect = thumbnail.getBoundingClientRect();
                return { x: rect.left, y: rect.top + pageYScroll, w: rect.width };
            }
        }
        let pswp = new PhotoSwipe(document.querySelectorAll('.pswp')[0],
            PhotoSwipeUI_Default, gallery, options);
        pswp.init();
    };
    // 图片点击打开 gallry
    let figureClick = function (e) {
        e = e || window.event;
        e.preventDefault ? e.preventDefault() : e.returnValue = false;
        let eTarget = e.target || e.srcElement;
        // 找到最近的 item
        let figure = eTarget;
        while (!figure.matches(itemSelector)) {
            figure = figure.parentNode;
        }
        // 找到最近的 gallery
        let gallery = figure;
        while (!gallery.matches(gallerySelector)) {
            gallery = gallery.parentNode;
        }
        let galleryIndex = parseInt(gallery.dataset.pswpGid) - 1;
        let itemIndex = parseInt(figure.dataset.pswpPid) - 1;
        openPhotoSwipe(galleryIndex, itemIndex);
    }
    // 初始化 photoSwipe
    let initPhotoSwipeFromDOM = function () {
        this.gallerys = [];
        // 设置宽度比例
        let photoList = document.querySelectorAll('.photos');
        for (let i = 0; i < photoList.length; i ++) {
            let scales = 0;
            let imgs = photoList[i].querySelectorAll('img');
            for (let j = 0; j < imgs.length; j ++) {
                scales += parseInt(imgs[j].getAttribute('width'));
            }
            for (let j = 0; j < imgs.length; j ++) {
                const x = parseInt(imgs[j].getAttribute('width')) / parseInt(scales) * 100;
                imgs[j].parentElement.style.flex = x;
            }
        }
        // 设置图片序号
        let gallyers = document.querySelectorAll(gallerySelector);
        for (let i = 0; i < gallyers.length; i ++) {
            let gallery = gallyers[i];
            gallery.setAttribute('data-pswp-gid', i + 1);
            let currentList = [];
            let figures = gallery.querySelectorAll(itemSelector);
            for (let j = 0; j < figures.length; j ++) {
                var figure = figures[j];
                figure.setAttribute('data-pswp-pid', j + 1);
                figure.onclick = figureClick;
                img = figure.firstElementChild;
                currentList.push({
                    src: img.src,
                    h: parseInt(img.getAttribute('height')) || img.height,
                    w: parseInt(img.getAttribute('width')) || img.width,
                    el: img
                })
            }
            this.gallerys.push(currentList);
        }
    }
    initPhotoSwipeFromDOM();
})('.pswp-gallery', '.pswp-item')