document.write('<div id="translate_main"></div>');

var clientKey = "20140727183457";

var srcLangList = jQuery('<select id="translate_src_lang"></select>');
jQuery(srcLangList).append(jQuery('<option value="en">英语</option>'));
jQuery(srcLangList).append(jQuery('<option value="fr">法语</option>'));
jQuery(srcLangList).append(jQuery('<option value="de">德语</option>'));
jQuery(srcLangList).append(jQuery('<option value="es">西班牙语</option>'));

jQuery(srcLangList).append(jQuery('<option value="zh-CN">简体中文</option>'));
jQuery(srcLangList).append(jQuery('<option value="zh-TW">繁体中文</option>'));

var destLangList = jQuery('<select id="translate_dest_lang"></select>');
jQuery(destLangList).append(jQuery('<option value="zh-CN">简体中文</option>'));
jQuery(destLangList).append(jQuery('<option value="zh-TW">繁体中文</option>'));

jQuery(destLangList).append(jQuery('<option value="en">英语</option>'));
jQuery(destLangList).append(jQuery('<option value="fr">法语</option>'));
jQuery(destLangList).append(jQuery('<option value="de">德语</option>'));
jQuery(destLangList).append(jQuery('<option value="es">西班牙语</option>'));

var srcText = jQuery('<textarea id="translate_src" rows="4" cols="80"></textarea>');
var destText = jQuery('<div id="translate_dest"></div>');

var translateButton = jQuery('<button>翻译</button>');
jQuery(translateButton).click(function() {
  var url = location.protocol + "//leonax.net/gservice/translate";
  url += "?q=" + encodeURIComponent(jQuery(srcText).val());
  url += "&s=" + jQuery(srcLangList).val();
  url += "&t=" + jQuery(destLangList).val();
  url += "&ck=" + clientKey;
  jQuery.post(url, function(result) {
    var text = result.data.translations[0].translatedText;
    text = text.replace("\n", "<br/>")
    destText.html(text);
  });
});

var srcLangOptionPanel = jQuery('<span></span>').append('<span>源语言：</span>').append(srcLangList);
var destLangOptionPanel = jQuery('<span></span>').append('<span>目标语言：</span>').append(destLangList);

var optionPanel = jQuery('<div></div>');
jQuery(optionPanel)
    .append(srcLangOptionPanel)
    .append('<span style="margin-right:50px;"></span>')
    .append(destLangOptionPanel)
    .append('<span style="margin-right:50px;"></span>')
    .append(translateButton);

var destOptionPanel = jQuery('<div>翻译结果：</div>');

jQuery('#translate_main')
    .append(optionPanel)
    .append(srcText)
    .append(destOptionPanel)
    .append(destText);
