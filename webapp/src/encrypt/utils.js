
const base64Tester = RegExp('^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$');

export function validBase64(str) {
    return base64Tester.test(str);
}
