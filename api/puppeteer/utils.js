const  toJSON = (node) => {
    node = node || this;

    const obj = {
        nodeType: node.nodeType,
        innerText: node.innerText
    };

    if (node.tagName) {
        obj.tagName = node.tagName.toLowerCase();
    }

    if (node.nodeName) {
        obj.nodeName = node.nodeName;
    }

    if (node.nodeValue) {
        obj.nodeValue = node.nodeValue;
    }

    const attrs = node.attributes;
    if (attrs) {
        let length = attrs.length;
        const arr = obj.attributes = new Array(length);
        for (var i = 0; i < length; i++) {
            attr = attrs[i];
            arr[i] = [attr.nodeName, attr.nodeValue];
        }
    }

    return obj;
}

module.exports = {
    toJSON
}
