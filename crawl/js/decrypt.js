const CryptoJS = require("./crawl/js/crypto-js.js");

function _0x442c(_0x1c5d27, _0x1eb2df) {
    const _0x2f5390 = _0x2f53();
    return _0x442c = function (_0x442cb3, _0x2e6e94) {
        _0x442cb3 = _0x442cb3 - 0x16d;
        let _0x1db293 = _0x2f5390[_0x442cb3];
        return _0x1db293;
    }, _0x442c(_0x1c5d27, _0x1eb2df);
}

(function (_0x32fd6e, _0x4692e5) {
    const _0x1eceb2 = _0x442c, _0x1c6ed7 = _0x32fd6e();
    while (!![]) {
        try {
            const _0x421b37 = parseInt(_0x1eceb2(0x175)) / 0x1 * (-parseInt(_0x1eceb2(0x170)) / 0x2) + -parseInt(_0x1eceb2(0x17c)) / 0x3 * (parseInt(_0x1eceb2(0x184)) / 0x4) + -parseInt(_0x1eceb2(0x171)) / 0x5 + -parseInt(_0x1eceb2(0x180)) / 0x6 * (-parseInt(_0x1eceb2(0x16e)) / 0x7) + parseInt(_0x1eceb2(0x17d)) / 0x8 * (parseInt(_0x1eceb2(0x17f)) / 0x9) + parseInt(_0x1eceb2(0x183)) / 0xa + -parseInt(_0x1eceb2(0x181)) / 0xb * (-parseInt(_0x1eceb2(0x172)) / 0xc);
            if (_0x421b37 === _0x4692e5) break; else _0x1c6ed7['push'](_0x1c6ed7['shift']());
        } catch (_0x407f5a) {
            _0x1c6ed7['push'](_0x1c6ed7['shift']());
        }
    }
}(_0x2f53, 0x5ebdf));

function decryptImage(_0x11ceb3) {
    const _0x4ce878 = _0x442c;
    try {
        const _0x19e68b = CryptoJS[String['fromCharCode'](0x65) + String[_0x4ce878(0x182)](0x6e) + String['fromCharCode'](0x63)][String[_0x4ce878(0x182)](0x55) + String['fromCharCode'](0x74) + String[_0x4ce878(0x182)](0x66) + String[_0x4ce878(0x182)](0x38)][String[_0x4ce878(0x182)](0x70) + _0x4ce878(0x179)](_0x4ce878(0x176)[_0x4ce878(0x174)]('_')[_0x4ce878(0x16d)](_0x11164b => String['fromCharCode'](parseInt(_0x11164b)))[_0x4ce878(0x178)]('')),
            _0x2d2822 = CryptoJS[String[_0x4ce878(0x182)](0x65) + String['fromCharCode'](0x6e) + String[_0x4ce878(0x182)](0x63)][String['fromCharCode'](0x55) + String[_0x4ce878(0x182)](0x74) + String['fromCharCode'](0x66) + String[_0x4ce878(0x182)](0x38)][String[_0x4ce878(0x182)](0x70) + _0x4ce878(0x179)](_0x4ce878(0x16f)[_0x4ce878(0x174)]('_')[_0x4ce878(0x16d)](_0x50224b => String['fromCharCode'](parseInt(_0x50224b)))[_0x4ce878(0x178)]('')),
            _0x1ef406 = CryptoJS['AES']['decrypt'](_0x11ceb3, _0x19e68b, {
                'iv': _0x2d2822,
                'mode': CryptoJS[_0x4ce878(0x17b)]['CBC'],
                'padding': CryptoJS[_0x4ce878(0x173)]['Pkcs7']
            });
        return _0x1ef406[_0x4ce878(0x17e)](CryptoJS[_0x4ce878(0x17a)][_0x4ce878(0x177)]);
    } catch (_0x544ab2) {
        console.log(_0x544ab2)
        return '';
    }
}

function _0x2f53() {
    const _0x38ae7a = ['Base64', 'join', 'arse', 'enc', 'mode', '1119fpkgWP', '16geOHAu', 'toString', '2758293xEIARY', '6peaOse', '6075773QCsjiV', 'fromCharCode', '3914210bXiumy', '8144yFccIN', 'map', '2827027mINBQi', '57_55_98_54_48_51_57_52_97_98_99_50_102_98_101_49', '4kGLuyr', '1304260BLiopi', '12NdebZO', 'pad', 'split', '276118tUSBDS', '102_53_100_57_54_53_100_102_55_53_51_51_54_50_55_48'];
    _0x2f53 = function () {
        return _0x38ae7a;
    };
    return _0x2f53();
};

// 从命令行参数中获取输入
const input = process.argv[2];

// 使用 CryptoJS 对输入进行加密
const img = decryptImage(input);

// 将加密后的文本输出
console.log(img);
