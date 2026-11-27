'use strict';
// Seller Center: produk — daftar, tambah/edit (varian, video/foto, kondisi), terjemahan AI.
addStrings({
  en: {wz_media: 'Video & photos', wz_basic: 'Name & category', wz_price: 'Price & stock', wz_check: 'Check & send', wzStepOf: 'Step {a} of 4', wzNext: 'Next: {s}', wzBack: 'Back', wzEdit: 'Change',
    wzh_media: 'Record the product in your hand, turn it around and show every colour. A short video helps buyers trust it — clear photos work too.', wzh_basic: 'Write the name the way buyers search for it, then pick a category.', wzh_price: 'Set the price buyers pay and how many you have. We show what you receive after the Phsar fee.', wzh_check: 'This is how buyers will see it. Check, then publish.',
    gv_title: 'A good video:', gv1: 'Bright light, plain background', gv2: 'Show the size next to your hand', gv3: 'Show every colour you sell',
    yg_title: 'You receive', yg_from: 'You receive from', yg_line: 'Price ${p} − Phsar fee {r}% (${f})', yg_none: 'Type a price to see what you receive.',
    rv_media: 'Video and photos', rv_basic: 'Name and category', rv_price: 'Price and stock', rv_ok: 'Ready', rv_bad: 'Missing', tabProducts: 'Products', addProduct: 'Add product', editProduct: 'Edit product', noProducts: 'No products here yet.',
    f_all: 'All', f_draft: 'Drafts', f_pending: 'In review', f_live: 'Live', f_changes: 'Needs changes', f_archived: 'Archived',
    ps_draft: 'Draft', ps_review: 'In review', ps_live: 'Live', ps_changes: 'Needs changes', ps_rejected: 'Rejected', ps_archived: 'Archived', ps_banned: 'Blocked', ps_hidden: 'Hidden (store not active)',
    stock: 'Stock', variantsN: '{n} variants', backProducts: '← Products', storeNotActiveHint: 'You can prepare products as drafts now. They can be sent for review once your store is approved.',
    secMedia: 'Video & photos', video: 'Product video', videoHint: 'Recommended — not required. Up to 3 minutes, max 100 MB. Products with a video get the “Real video” badge and appear in the Videos feed.',
    uploadVideo: 'Upload video', replace: 'Replace', remove: 'Remove', photos: 'Photos', photosHint: 'Up to 9. The first photo is the cover. No video? Add at least 3 photos from different angles.', addPhotos: 'Add photos',
    defects: 'Condition photos', defectsHint: 'Close-ups of scratches, wear or missing parts. Required for used items.', moveFirst: 'Make cover',
    uploading: 'Uploading… {p}%', videoTooLong: 'The video is longer than 3 minutes.', fileTooBig: 'The file is too large.',
    secBasic: 'Product details', productName: 'Product name', namePh: 'e.g. Cotton T-shirt, oversized', category: 'Category', chooseCategory: 'Choose a subcategory…', catHint: 'Pick the most specific subcategory — buyers filter by it.', needSubcat: 'Choose a subcategory for this product (categories now have subcategories).', description: 'Description',
    descPh: 'Material, size, what’s in the box…', translateHint: 'Write in one language only. We translate the name, description, colours and sizes into the other two languages automatically.',
    secCondition: 'Condition', condition: 'Condition', c_new: 'New', c_like_new: 'Like new', c_used_good: 'Used – good', c_used_fair: 'Used – fair', c_refurbished: 'Refurbished',
    conditionNote: 'Describe the condition honestly', conditionNotePh: 'e.g. small scratch on the back, box missing', usageMonths: 'Used for (months)', newOnlyCategory: 'This category is for new items only.',
    secVariants: 'Price, stock & variants', hasOptions: 'This product comes in different colours, sizes…', price: 'Price (USD)', comparePrice: 'Original price', stockQty: 'Stock', sku: 'SKU',
    optional: 'optional', option: 'Option', optKind: 'Type', k_color: 'Colour', k_size: 'Size', k_material: 'Material', k_capacity: 'Capacity', k_custom: 'Other',
    optName: 'Name shown to buyers', addValue: 'Add', valuePh: 'Type and press Enter', addOption: 'Add another option', removeOption: 'Remove option',
    available: 'Sell', applyAll: 'Apply to all', variant: 'Variant', noVariants: 'Add at least one value to each option.',
    secMore: 'Warranty & returns', warranty: 'Warranty', w_none: 'No warranty', w_store: 'Store warranty', w_official: 'Official warranty', warrantyMonths: 'Months', returnDays: 'Return within (days)',
    saveDraft: 'Save draft', submitReview: 'Publish', saveResubmit: 'Save', reviewNote: 'Every change is checked automatically. Price and stock changes go live right away; changing the name, category, condition or cover photo is re-checked in full. A very low price or an unusual discount may wait for Phsar to check it.',
    savedDraft: 'Draft saved', sentReview: 'Sent for review', archive: 'Hide from shop', unarchive: 'Show again', confirmDelete: 'Delete this product?',
    modPending: 'Phsar is checking this product (usually within a few hours).', modChanges: 'Phsar asked for changes:', modRejected: 'This product can’t be sold on Phsar:', modLive: 'Live — buyers can see this product.', nowLive: 'Published — it’s live now', modWhy: 'Why it’s checked by hand:', priceHuge: 'The price is {p}. Is that right?', srs_new_store: 'your store’s first products are checked by Phsar — after that, products go live straight away', srs_restricted_category: 'this category is always checked by hand', srs_keyword: 'the name or description contains a word we check by hand', srs_price_floor: 'the price looks unusually low for this item', srs_big_discount: 'the original (crossed-out) price is much higher than the price', srs_price_drop: 'the price dropped a lot since it was approved', srs_store_strikes: 'your store has active strikes', srs_changes_requested: 'Phsar asked for changes — we check your fix', srs_was_rejected: 'this product was rejected before', srs_reported: 'buyers reported this product', srs_auto_publish_off: 'all products are being checked by hand right now',
    secTranslations: 'Translations', translating: 'Translating…', tr_machine: 'AI translation', tr_needs_review: 'AI — please check', tr_seller: 'Edited by you', tr_reviewer: 'Checked by Phsar', tr_source: 'Your original',
    editTr: 'Edit', useAi: 'Use AI translation', translationsLater: 'Translations appear here after you save.',
    err_store_not_active: 'Your store isn’t active yet. Save as draft for now.',
    compressing: 'Compressing video…', compressed: 'optimised', optimizing: 'Optimising…', trs_pending: 'Translating…', trs_retrying: 'Translation problem', trs_failed: 'Translation failed', trs_off: 'Auto-translate off', trOffShort: 'Not translated', trFailedShort: 'Failed', writeTr: 'Write translation', trOff: 'Automatic translation is not set up on this server (ANTHROPIC_API_KEY is empty). You can write the translations yourself, or ask the Phsar team.', trFailed: 'The automatic translation failed. Try again, or write the translation yourself.', trRetrying: 'The translation service had a problem. We’ll retry automatically in a few minutes.', trRetry: 'Try again now', submitLocked: 'Send for review unlocks after Phsar approves your store.', draftVerify: 'Save draft & verify my store', draftVerifyHint: 'Your product is saved as a draft and you go straight to the store check. Once Phsar approves your store, you can publish it with one tap.', draftVerifyFix: 'Draft saved. Finish this step so Phsar can verify your store.', storeReviewing: 'Your store is being reviewed by Phsar. You can publish as soon as it is approved.', compareHigher: 'Original price must be higher than the price (it’s shown crossed out). Leave it empty if there’s no discount.', needVideo: 'Add one product video.', needPhoto: 'Add at least one photo.', needMorePhotos: 'Without a video, add at least 3 photos (front, back, details).', noVideoNudge: 'No video yet — that’s okay. A 15–30 second phone clip helps buyers trust the item and puts it in the Videos feed.', vidOptional: 'optional', noVideoChip: 'No video', photoTips: 'Good photos: ✓ bright daylight ✓ front, back & sides ✓ close-up of labels or details', needDefect: 'Used items need a close-up photo of any scratches or wear.', needValue: 'Add at least one value, e.g. Red or M.', needVariantOn: 'Tick at least one variant to sell.', err_product_banned: 'This product was blocked by Phsar.'},
  km: {wz_media: 'វីដេអូ និងរូបថត', wz_basic: 'ឈ្មោះ និងប្រភេទ', wz_price: 'តម្លៃ និងស្តុក', wz_check: 'ពិនិត្យ និងផ្ញើ', wzStepOf: 'ជំហាន {a} នៃ 4', wzNext: 'បន្ទាប់៖ {s}', wzBack: 'ថយក្រោយ', wzEdit: 'ប្តូរ',
    wzh_media: 'ថតផលិតផលនៅក្នុងដៃ បង្វិលវា ហើយបង្ហាញគ្រប់ពណ៌។ វីដេអូខ្លីជួយឲ្យអ្នកទិញទុកចិត្ត — រូបថតច្បាស់ក៏បានដែរ។', wzh_basic: 'សរសេរឈ្មោះតាមរបៀបដែលអ្នកទិញស្វែងរក រួចជ្រើសប្រភេទ។', wzh_price: 'កំណត់តម្លៃអ្នកទិញបង់ និងចំនួនដែលអ្នកមាន។ យើងបង្ហាញចំនួនអ្នកទទួលបានបន្ទាប់ពីកម្រៃ Phsar។', wzh_check: 'នេះជារបៀបដែលអ្នកទិញនឹងឃើញ។ ពិនិត្យ រួចដាក់លក់។',
    gv_title: 'វីដេអូល្អ៖', gv1: 'ពន្លឺភ្លឺ ផ្ទៃខាងក្រោយសាមញ្ញ', gv2: 'បង្ហាញទំហំក្បែរដៃ', gv3: 'បង្ហាញគ្រប់ពណ៌ដែលលក់',
    yg_title: 'អ្នកទទួលបាន', yg_from: 'អ្នកទទួលបានចាប់ពី', yg_line: 'តម្លៃ ${p} − កម្រៃ Phsar {r}% (${f})', yg_none: 'វាយតម្លៃ ដើម្បីមើលចំនួនដែលអ្នកទទួលបាន។',
    rv_media: 'វីដេអូ និងរូបថត', rv_basic: 'ឈ្មោះ និងប្រភេទ', rv_price: 'តម្លៃ និងស្តុក', rv_ok: 'រួចរាល់', rv_bad: 'ខ្វះ', tabProducts: 'ផលិតផល', addProduct: 'បន្ថែមផលិតផល', editProduct: 'កែផលិតផល', noProducts: 'មិនទាន់មានផលិតផលទេ។',
    f_all: 'ទាំងអស់', f_draft: 'សេចក្តីព្រាង', f_pending: 'កំពុងពិនិត្យ', f_live: 'កំពុងលក់', f_changes: 'ត្រូវកែ', f_archived: 'បានលាក់',
    ps_draft: 'សេចក្តីព្រាង', ps_review: 'កំពុងពិនិត្យ', ps_live: 'កំពុងលក់', ps_changes: 'ត្រូវកែ', ps_rejected: 'បានបដិសេធ', ps_archived: 'បានលាក់', ps_banned: 'បានហាមឃាត់', ps_hidden: 'លាក់ (ហាងមិនទាន់សកម្ម)',
    stock: 'ស្តុក', variantsN: '{n} ជម្រើស', backProducts: '← ផលិតផល', storeNotActiveHint: 'អ្នកអាចរៀបចំផលិតផលជាសេចក្តីព្រាងសិន។ អាចផ្ញើឱ្យពិនិត្យបាន ពេលហាងត្រូវបានអនុម័ត។',
    secMedia: 'វីដេអូ និងរូបថត', video: 'វីដេអូផលិតផល', videoHint: 'ណែនាំ — មិនចាំបាច់។ យូរបំផុត 3 នាទី ទំហំ 100 MB។ ផលិតផលមានវីដេអូ ទទួលបានស្លាក «វីដេអូពិត» និងបង្ហាញក្នុងផ្នែកវីដេអូ។',
    uploadVideo: 'ផ្ទុកវីដេអូ', replace: 'ប្តូរ', remove: 'ដកចេញ', photos: 'រូបថត', photosHint: 'ច្រើនបំផុត 9។ រូបទីមួយជារូបគម្រប។ គ្មានវីដេអូ? សូមបន្ថែមរូបយ៉ាងហោចណាស់ 3 ពីមុំផ្សេងៗ។', addPhotos: 'បន្ថែមរូបថត',
    defects: 'រូបថតស្ថានភាព', defectsHint: 'ថតជិតកន្លែងឆ្កូត ប្រេះ ឬខ្វះ។ ចាំបាច់សម្រាប់ទំនិញប្រើរួច។', moveFirst: 'ធ្វើជារូបគម្រប',
    uploading: 'កំពុងផ្ទុក… {p}%', videoTooLong: 'វីដេអូលើសពី 3 នាទី។', fileTooBig: 'ឯកសារធំពេក។',
    secBasic: 'ព័ត៌មានផលិតផល', productName: 'ឈ្មោះផលិតផល', namePh: 'ឧ. អាវយឺតកប្បាស', category: 'ប្រភេទ', chooseCategory: 'ជ្រើសរើសប្រភេទរង…', catHint: 'ជ្រើសប្រភេទរងដែលត្រឹមត្រូវបំផុត — អ្នកទិញច្រោះតាមវា។', needSubcat: 'សូមជ្រើសប្រភេទរងសម្រាប់ផលិតផលនេះ (ប្រភេទឥឡូវមានប្រភេទរង)។', description: 'ការពិពណ៌នា',
    descPh: 'សាច់ក្រណាត់ ទំហំ អ្វីខ្លះក្នុងប្រអប់…', translateHint: 'សរសេរតែមួយភាសា។ យើងបកប្រែឈ្មោះ ការពិពណ៌នា ពណ៌ និងទំហំ ទៅភាសាពីរទៀតដោយស្វ័យប្រវត្តិ។',
    secCondition: 'ស្ថានភាព', condition: 'ស្ថានភាព', c_new: 'ថ្មី', c_like_new: 'ដូចថ្មី', c_used_good: 'ប្រើរួច – ល្អ', c_used_fair: 'ប្រើរួច – មធ្យម', c_refurbished: 'ជួសជុលឡើងវិញ',
    conditionNote: 'ពិពណ៌នាស្ថានភាពដោយស្មោះត្រង់', conditionNotePh: 'ឧ. ឆ្កូតតិចនៅខាងក្រោយ គ្មានប្រអប់', usageMonths: 'បានប្រើ (ខែ)', newOnlyCategory: 'ប្រភេទនេះសម្រាប់តែទំនិញថ្មី។',
    secVariants: 'តម្លៃ ស្តុក និងជម្រើស', hasOptions: 'ផលិតផលនេះមានពណ៌ ទំហំ… ផ្សេងៗ', price: 'តម្លៃ (USD)', comparePrice: 'តម្លៃដើម', stockQty: 'ស្តុក', sku: 'SKU',
    optional: 'មិនចាំបាច់', option: 'ជម្រើស', optKind: 'ប្រភេទ', k_color: 'ពណ៌', k_size: 'ទំហំ', k_material: 'សាច់', k_capacity: 'ចំណុះ', k_custom: 'ផ្សេងទៀត',
    optName: 'ឈ្មោះបង្ហាញអ្នកទិញ', addValue: 'បន្ថែម', valuePh: 'វាយ ហើយចុច Enter', addOption: 'បន្ថែមជម្រើសមួយទៀត', removeOption: 'ដកជម្រើស',
    available: 'លក់', applyAll: 'ដាក់ទាំងអស់', variant: 'ជម្រើស', noVariants: 'បន្ថែមតម្លៃយ៉ាងហោចមួយក្នុងជម្រើសនីមួយៗ។',
    secMore: 'ធានា និងការប្តូរ', warranty: 'ការធានា', w_none: 'គ្មានការធានា', w_store: 'ធានាដោយហាង', w_official: 'ធានាផ្លូវការ', warrantyMonths: 'ខែ', returnDays: 'ប្តូរក្នុងរយៈ (ថ្ងៃ)',
    saveDraft: 'រក្សាទុកសេចក្តីព្រាង', submitReview: 'ដាក់លក់', saveResubmit: 'រក្សាទុក', reviewNote: 'រាល់ការកែប្រែត្រូវពិនិត្យដោយស្វ័យប្រវត្តិ។ ការកែតម្លៃ និងស្តុកមានប្រសិទ្ធភាពភ្លាមៗ ការកែឈ្មោះ ប្រភេទ ស្ថានភាព ឬរូបគម្រប ត្រូវពិនិត្យឡើងវិញទាំងស្រុង។ តម្លៃទាបខ្លាំង ឬការបញ្ចុះតម្លៃមិនធម្មតា អាចរង់ចាំ Phsar ពិនិត្យ។',
    savedDraft: 'បានរក្សាទុកសេចក្តីព្រាង', sentReview: 'បានផ្ញើឱ្យពិនិត្យ', archive: 'លាក់ពីហាង', unarchive: 'បង្ហាញវិញ', confirmDelete: 'លុបផលិតផលនេះ?',
    modPending: 'Phsar កំពុងពិនិត្យផលិតផលនេះ (ជាធម្មតាក្នុងពីរបីម៉ោង)។', modChanges: 'Phsar ស្នើឱ្យកែ៖', modRejected: 'ផលិតផលនេះមិនអាចលក់លើ Phsar បានទេ៖', modLive: 'កំពុងលក់ — អ្នកទិញមើលឃើញផលិតផលនេះ។', nowLive: 'បានបង្ហាញ — កំពុងលក់ឥឡូវនេះ', modWhy: 'ហេតុអ្វីត្រូវពិនិត្យដោយដៃ៖', priceHuge: 'តម្លៃគឺ {p}។ ត្រឹមត្រូវទេ?', srs_new_store: 'ផលិតផលដំបូងៗរបស់ហាងត្រូវ Phsar ពិនិត្យ — បន្ទាប់មក ផលិតផលនឹងបង្ហាញភ្លាមៗ', srs_restricted_category: 'ប្រភេទនេះតែងតែពិនិត្យដោយដៃ', srs_keyword: 'ឈ្មោះ ឬការពិពណ៌នាមានពាក្យដែលយើងពិនិត្យដោយដៃ', srs_price_floor: 'តម្លៃហាក់ទាបខុសធម្មតាសម្រាប់ទំនិញនេះ', srs_big_discount: 'តម្លៃដើម (គូសចោល) ខ្ពស់ជាងតម្លៃលក់ច្រើនណាស់', srs_price_drop: 'តម្លៃធ្លាក់ខ្លាំងចាប់តាំងពីអនុម័ត', srs_store_strikes: 'ហាងរបស់អ្នកមានការព្រមានសកម្ម', srs_changes_requested: 'Phsar បានស្នើឲ្យកែ — យើងពិនិត្យការកែរបស់អ្នក', srs_was_rejected: 'ផលិតផលនេះធ្លាប់ត្រូវបដិសេធ', srs_reported: 'អ្នកទិញបានរាយការណ៍ផលិតផលនេះ', srs_auto_publish_off: 'ឥឡូវនេះ ផលិតផលទាំងអស់ត្រូវពិនិត្យដោយដៃ',
    secTranslations: 'ការបកប្រែ', translating: 'កំពុងបកប្រែ…', tr_machine: 'បកប្រែដោយ AI', tr_needs_review: 'AI — សូមពិនិត្យ', tr_seller: 'អ្នកបានកែ', tr_reviewer: 'Phsar បានពិនិត្យ', tr_source: 'អត្ថបទដើមរបស់អ្នក',
    editTr: 'កែ', useAi: 'ប្រើការបកប្រែ AI', translationsLater: 'ការបកប្រែនឹងបង្ហាញនៅទីនេះ បន្ទាប់ពីរក្សាទុក។',
    err_store_not_active: 'ហាងមិនទាន់សកម្ម។ សូមរក្សាទុកជាសេចក្តីព្រាងសិន។',
    compressing: 'កំពុងបង្រួមវីដេអូ…', compressed: 'បានបង្រួម', optimizing: 'កំពុងរៀបចំ…', trs_pending: 'កំពុងបកប្រែ…', trs_retrying: 'បញ្ហាបកប្រែ', trs_failed: 'បកប្រែមិនបាន', trs_off: 'ការបកប្រែស្វ័យប្រវត្តិបិទ', trOffShort: 'មិនទាន់បកប្រែ', trFailedShort: 'បរាជ័យ', writeTr: 'សរសេរការបកប្រែ', trOff: 'ម៉ាស៊ីនមេនេះមិនទាន់កំណត់ការបកប្រែស្វ័យប្រវត្តិ។ អ្នកអាចសរសេរការបកប្រែដោយខ្លួនឯង ឬសួរក្រុម Phsar។', trFailed: 'ការបកប្រែស្វ័យប្រវត្តិបរាជ័យ។ សូមព្យាយាមម្តងទៀត ឬសរសេរដោយខ្លួនឯង។', trRetrying: 'សេវាបកប្រែមានបញ្ហា។ យើងនឹងព្យាយាមម្តងទៀតក្នុងពីរបីនាទី។', trRetry: 'ព្យាយាមឥឡូវនេះ', submitLocked: 'អាចផ្ញើឱ្យពិនិត្យបាន ពេល Phsar អនុម័តហាងរបស់អ្នក។', draftVerify: 'រក្សាទុកសេចក្តីព្រាង និងផ្ទៀងផ្ទាត់ហាង', draftVerifyHint: 'ផលិតផលត្រូវបានរក្សាទុកជាសេចក្តីព្រាង ហើយអ្នកទៅផ្ទៀងផ្ទាត់ហាងភ្លាម។ ពេល Phsar អនុម័តហាង អ្នកអាចផ្សាយវាដោយចុចម្តង។', draftVerifyFix: 'បានរក្សាទុកសេចក្តីព្រាង។ សូមបំពេញជំហាននេះ ដើម្បីឱ្យ Phsar ផ្ទៀងផ្ទាត់ហាងរបស់អ្នក។', storeReviewing: 'Phsar កំពុងពិនិត្យហាងរបស់អ្នក។ អ្នកអាចផ្សាយបានភ្លាមៗពេលត្រូវបានអនុម័ត។', compareHigher: 'តម្លៃដើមត្រូវខ្ពស់ជាងតម្លៃលក់ (បង្ហាញជាតម្លៃគូសចោល)។ ទុកទទេ បើគ្មានការបញ្ចុះតម្លៃ។', needVideo: 'សូមបន្ថែមវីដេអូផលិតផលមួយ។', needPhoto: 'សូមបន្ថែមរូបថតយ៉ាងហោចណាស់មួយ។', needMorePhotos: 'បើគ្មានវីដេអូ សូមបន្ថែមរូបយ៉ាងហោចណាស់ 3 (មុខ ក្រោយ ព័ត៌មានលម្អិត)។', noVideoNudge: 'មិនទាន់មានវីដេអូ — មិនអីទេ។ ឃ្លីបទូរសព្ទ 15–30 វិនាទី ជួយឲ្យអ្នកទិញទុកចិត្ត ហើយបង្ហាញក្នុងផ្នែកវីដេអូ។', vidOptional: 'ស្រេចចិត្ត', noVideoChip: 'គ្មានវីដេអូ', photoTips: 'រូបល្អ៖ ✓ ពន្លឺថ្ងៃភ្លឺ ✓ មុខ ក្រោយ និងចំហៀង ✓ ថតជិតស្លាក ឬព័ត៌មានលម្អិត', needDefect: 'ទំនិញប្រើរួចត្រូវមានរូបថតជិតកន្លែងឆ្កូត ឬខូច។', needValue: 'សូមបន្ថែមយ៉ាងហោចមួយ ឧ. ក្រហម ឬ M។', needVariantOn: 'សូមធីកជម្រើសយ៉ាងហោចមួយដើម្បីលក់។', err_product_banned: 'ផលិតផលនេះត្រូវបាន Phsar ហាមឃាត់។'},
  zh: {wz_media: '视频和照片', wz_basic: '名称和分类', wz_price: '价格和库存', wz_check: '检查并提交', wzStepOf: '第 {a} 步，共 4 步', wzNext: '下一步：{s}', wzBack: '返回', wzEdit: '修改',
    wzh_media: '把商品拿在手里拍，转一圈，展示每种颜色。短视频更能让买家放心——清晰的照片也可以。', wzh_basic: '按买家搜索的方式写名称，然后选择分类。', wzh_price: '设置买家支付的价格和库存数量。我们会显示扣除 Phsar 佣金后你能收到多少。', wzh_check: '这是买家看到的样子。检查后发布。',
    gv_title: '好视频：', gv1: '光线明亮、背景简洁', gv2: '放在手边展示尺寸', gv3: '展示所有在售颜色',
    yg_title: '你将收到', yg_from: '你将收到（起）', yg_line: '价格 ${p} − Phsar 佣金 {r}%（${f}）', yg_none: '输入价格即可查看你能收到多少。',
    rv_media: '视频和照片', rv_basic: '名称和分类', rv_price: '价格和库存', rv_ok: '已完成', rv_bad: '缺少', tabProducts: '商品', addProduct: '添加商品', editProduct: '编辑商品', noProducts: '这里还没有商品。',
    f_all: '全部', f_draft: '草稿', f_pending: '审核中', f_live: '在售', f_changes: '需修改', f_archived: '已下架',
    ps_draft: '草稿', ps_review: '审核中', ps_live: '在售', ps_changes: '需修改', ps_rejected: '已拒绝', ps_archived: '已下架', ps_banned: '已封禁', ps_hidden: '未展示（店铺未上线）',
    stock: '库存', variantsN: '{n} 个规格', backProducts: '← 商品', storeNotActiveHint: '现在可以先保存为草稿，店铺审核通过后再提交审核。',
    secMedia: '视频和照片', video: '商品视频', videoHint: '建议上传，非必填。最长 3 分钟，最大 100 MB。有视频的商品会显示“实拍视频”标志，并进入视频频道。',
    uploadVideo: '上传视频', replace: '更换', remove: '删除', photos: '照片', photosHint: '最多 9 张，第一张为封面。没有视频？请至少上传 3 张不同角度的照片。', addPhotos: '添加照片',
    defects: '成色照片', defectsHint: '近拍划痕、磨损或缺件处。二手商品必填。', moveFirst: '设为封面',
    uploading: '上传中… {p}%', videoTooLong: '视频超过 3 分钟。', fileTooBig: '文件太大。',
    secBasic: '商品信息', productName: '商品名称', namePh: '例如：纯棉宽松T恤', category: '类目', chooseCategory: '请选择子类目…', catHint: '请选择最具体的子类目——买家会按它筛选。', needSubcat: '请为此商品选择子类目（类目现在包含子类目）。', description: '商品描述',
    descPh: '材质、尺寸、包装内容…', translateHint: '只需用一种语言填写。名称、描述、颜色和尺码会自动翻译成另外两种语言。',
    secCondition: '成色', condition: '成色', c_new: '全新', c_like_new: '几乎全新', c_used_good: '二手 – 良好', c_used_fair: '二手 – 一般', c_refurbished: '翻新',
    conditionNote: '如实描述成色', conditionNotePh: '例如：背面有轻微划痕，无包装盒', usageMonths: '已使用（月）', newOnlyCategory: '该类目只允许全新商品。',
    secVariants: '价格、库存和规格', hasOptions: '该商品有不同颜色、尺码…', price: '价格（美元）', comparePrice: '原价', stockQty: '库存', sku: 'SKU',
    optional: '选填', option: '规格', optKind: '类型', k_color: '颜色', k_size: '尺码', k_material: '材质', k_capacity: '容量', k_custom: '其他',
    optName: '向买家显示的名称', addValue: '添加', valuePh: '输入后按回车', addOption: '再加一个规格', removeOption: '删除规格',
    available: '在售', applyAll: '全部应用', variant: '规格', noVariants: '每个规格至少添加一个值。',
    secMore: '保修和退换', warranty: '保修', w_none: '无保修', w_store: '店铺保修', w_official: '官方保修', warrantyMonths: '月', returnDays: '可退换天数',
    saveDraft: '保存草稿', submitReview: '发布', saveResubmit: '保存', reviewNote: '所有修改都会自动检查。价格和库存修改立即生效；修改名称、类目、成色或封面照片会重新完整检查。价格过低或折扣异常可能需要等待 Phsar 审核。',
    savedDraft: '草稿已保存', sentReview: '已提交审核', archive: '从店铺下架', unarchive: '重新上架', confirmDelete: '删除此商品？',
    modPending: 'Phsar 正在审核此商品（通常几小时内）。', modChanges: 'Phsar 要求修改：', modRejected: '此商品不能在 Phsar 销售：', modLive: '在售——买家可以看到此商品。', nowLive: '已发布——现在已上架', modWhy: '需要人工审核的原因：', priceHuge: '价格是 {p}，确定吗？', srs_new_store: '店铺的前几件商品由 Phsar 审核——之后商品将直接上架', srs_restricted_category: '此类目始终人工审核', srs_keyword: '名称或描述包含需要人工审核的词语', srs_price_floor: '对此商品而言价格异常低', srs_big_discount: '原价（划线价）远高于售价', srs_price_drop: '价格较通过审核时大幅下降', srs_store_strikes: '你的店铺有生效中的违规记录', srs_changes_requested: 'Phsar 要求修改——我们会检查你的修改', srs_was_rejected: '此商品曾被拒绝', srs_reported: '有买家举报了此商品', srs_auto_publish_off: '目前所有商品都需人工审核',
    secTranslations: '翻译', translating: '翻译中…', tr_machine: 'AI 翻译', tr_needs_review: 'AI——请检查', tr_seller: '已由你修改', tr_reviewer: 'Phsar 已核对', tr_source: '你的原文',
    editTr: '编辑', useAi: '使用 AI 翻译', translationsLater: '保存后这里会显示翻译。',
    err_store_not_active: '店铺尚未上线，请先保存为草稿。',
    compressing: '正在压缩视频…', compressed: '已优化', optimizing: '处理中…', trs_pending: '翻译中…', trs_retrying: '翻译出现问题', trs_failed: '翻译失败', trs_off: '自动翻译未开启', trOffShort: '未翻译', trFailedShort: '失败', writeTr: '自己填写翻译', trOff: '此服务器未设置自动翻译（ANTHROPIC_API_KEY 为空）。你可以自己填写翻译，或联系 Phsar 团队。', trFailed: '自动翻译失败。请重试，或自己填写翻译。', trRetrying: '翻译服务出现问题，几分钟后会自动重试。', trRetry: '立即重试', submitLocked: 'Phsar 通过店铺审核后才能提交商品审核。', draftVerify: '保存草稿并验证店铺', draftVerifyHint: '商品会保存为草稿，然后直接进入店铺验证。Phsar 通过店铺审核后，一键即可发布。', draftVerifyFix: '草稿已保存。请完成这一步，以便 Phsar 验证你的店铺。', storeReviewing: 'Phsar 正在审核你的店铺，通过后即可发布。', compareHigher: '原价必须高于售价（会以划线价显示）。没有折扣请留空。', needVideo: '请添加一个商品视频。', needPhoto: '请至少添加一张照片。', needMorePhotos: '没有视频时，请至少添加 3 张照片（正面、背面、细节）。', noVideoNudge: '还没有视频——没关系。一段 15–30 秒的手机短片能让买家更放心，并进入视频频道。', vidOptional: '可选', noVideoChip: '无视频', photoTips: '好照片：✓ 明亮自然光 ✓ 正面、背面和侧面 ✓ 标签或细节特写', needDefect: '二手商品需要划痕或磨损处的近拍照片。', needValue: '请至少添加一个值，例如 红色 或 M。', needVariantOn: '请至少勾选一个要销售的规格。', err_product_banned: '此商品已被 Phsar 封禁。'},
});

const LANGS = {en: 'English', km: 'ខ្មែរ', zh: '中文'};
const money = v => v == null ? '—' : '$' + Number(v).toFixed(2);
const PF = {filter: 'all', cats: {}, E: null};

function productState(p, storeActive) {
  if (p.status === 'draft') return ['', 'ps_draft'];
  if (p.status === 'banned') return ['bad', 'ps_banned'];
  if (p.status === 'archived') return ['', 'ps_archived'];
  if (p.moderation_status === 'pending_review') return ['warn', 'ps_review'];
  if (p.moderation_status === 'changes_requested') return ['bad', 'ps_changes'];
  if (p.moderation_status === 'rejected') return ['bad', 'ps_rejected'];
  return storeActive ? ['ok', 'ps_live'] : ['warn', 'ps_hidden'];
}

// catIndex: id → {name, label, leaf, allows_used} untuk kategori utama & sub-kategori (bekas boleh = anak DAN induk).
function catIndex(cats) {
  const m = {};
  for (const c of cats || []) {
    const kids = c.children || [];
    m[c.id] = {id: c.id, name: c.name, label: c.name, leaf: !kids.length, allows_used: c.allows_used};
    for (const k of kids) m[k.id] = {id: k.id, name: k.name, label: `${c.name} › ${k.name}`, leaf: true, allows_used: c.allows_used && k.allows_used};
  }
  return m;
}
function catOptions(cats, cur) {
  const opt = c => `<option value="${c.id}" ${cur === c.id ? 'selected' : ''}>${esc(c.name)}</option>`;
  return (cats || []).map(c => (c.children || []).length ? `<optgroup label="${esc(c.name)}">${c.children.map(opt).join('')}</optgroup>` : opt(c)).join('');
}

async function categories() {
  if (!PF.cats[P.lang]) PF.cats[P.lang] = (await api('/v1/app/categories')).categories;
  return PF.cats[P.lang];
}

// ---- Tab daftar produk (di halaman toko) ----
addStrings({
  en: {pSearchPh: 'Search by product name or SKU', pSort: 'Sort', pSort_updated: 'Recently updated', pSort_name: 'Name A–Z', pSort_price_asc: 'Price: low to high',
    pSort_price_desc: 'Price: high to low', pSort_stock_low: 'Low stock first'},
  km: {pSearchPh: 'ស្វែងរកតាមឈ្មោះផលិតផល ឬ SKU', pSort: 'តម្រៀប', pSort_updated: 'ទើបកែប្រែ', pSort_name: 'ឈ្មោះ ក–អ', pSort_price_asc: 'តម្លៃ៖ ទាបទៅខ្ពស់',
    pSort_price_desc: 'តម្លៃ៖ ខ្ពស់ទៅទាប', pSort_stock_low: 'ស្តុកតិចមុន'},
  zh: {pSearchPh: '按商品名称或 SKU 搜索', pSort: '排序', pSort_updated: '最近更新', pSort_name: '名称 A–Z', pSort_price_asc: '价格从低到高',
    pSort_price_desc: '价格从高到低', pSort_stock_low: '库存少的优先'},
});
const PSORTS = ['updated', 'name', 'price_asc', 'price_desc', 'stock_low'];
Object.assign(PF, {q: '', sort: 'updated', page: 1});

async function tabProducts(id) {
  const query = qs({filter: PF.filter, q: PF.q, sort: PF.sort, page: PF.page, size: 24});
  const [list, cats] = await Promise.all([api(`/v1/seller/stores/${id}/products?${query}`), categories()]);
  const {products, counts = {}, page} = list;
  if (page && PF.page > page.pages) { PF.page = page.pages; return tabProducts(id); } // mis. setelah produk dihapus
  const CI = catIndex(cats), catName = c => (CI[c] || {}).name || c;
  const active = S.store.status === 'active';
  const cards = products.map(p => {
    const [cls, key] = productState(p, active);
    const price = p.min_price === p.max_price ? money(p.min_price) : `${money(p.min_price)} – ${money(p.max_price)}`;
    return `<div class="card pcard"><a class="plink" href="#/store/${id}/product/${p.id}">
      <div class="pcover">${p.cover_url ? `<img src="${esc(p.cover_url)}" alt="" loading="lazy">` : ''}${p.has_video ? '<span class="vbadge" aria-hidden="true">▶</span>' : p.cover_url ? `<span class="novid" title="${t('noVideoNudge')}">${t('noVideoChip')}</span>` : ''}</div>
      <div class="spread"><span class="chip ${cls}">${t(key)}</span><span class="muted">${esc(catName(p.category_id))}</span></div>
      <b class="pname">${esc(p.name)}</b>
      <div class="spread"><b>${price}</b><span class="muted">${t('stock')} ${p.total_stock}${p.variants > 1 ? ' · ' + t('variantsN', {n: p.variants}) : ''}</span></div>
      ${p.moderation_note && ['changes_requested', 'rejected'].includes(p.moderation_status) ? `<p class="muted" style="color:var(--bad)">${esc(p.moderation_note)}</p>` : ''}
    </a>${p.status !== 'banned' ? `<button type="button" class="small" onclick="openStockDialog('${id}','${p.id}')">${t('stockPrice')}</button>` : ''}</div>`;
  }).join('');
  const empty = PF.q
    ? `<section class="card"><p class="muted" style="text-align:center">${t('lsNoMatch', {q: esc(PF.q)})}</p></section>`
    : `<p class="muted">${t('noProducts')}</p>`;
  $('#tab').innerHTML = `
    ${active ? '' : `<div class="banner info">${t('storeNotActiveHint')}</div>`}
    <div class="spread"><div class="tabs">${['all', 'draft', 'pending', 'live', 'changes', 'archived'].map(f =>
      `<button type="button" aria-selected="${PF.filter === f}" data-pf="${f}">${t('f_' + f)}${counts[f] != null ? ` <span class="tcount">${counts[f]}</span>` : ''}</button>`).join('')}</div>
      <a class="btn primary" href="#/store/${id}/product/new">+ ${t('addProduct')}</a></div>
    <div class="listbar">${searchBox('pQ', PF.q, t('pSearchPh'))}
      <label class="sr-only" for="pSort">${t('pSort')}</label>
      <select id="pSort" aria-label="${t('pSort')}">${PSORTS.map(k => `<option value="${k}" ${PF.sort === k ? 'selected' : ''}>${t('pSort_' + k)}</option>`).join('')}</select></div>
    ${products.length ? `<div class="pgrid">${cards}</div>` : empty}
    ${pagerHTML(page)}`;
  document.querySelectorAll('[data-pf]').forEach(b => b.onclick = () => { PF.filter = b.dataset.pf; PF.page = 1; tabProducts(id); });
  $('#pSort').onchange = e => { PF.sort = e.target.value; PF.page = 1; tabProducts(id); };
  bindSearch('pQ', v => { if (v === PF.q) return; PF.q = v; PF.page = 1; tabProducts(id); });
  bindPager($('#tab'), p => { PF.page = p; tabProducts(id); });
}

// ---- Editor produk ----
function emptyProduct() {
  return {source_lang: P.lang, name: '', description: '', category_id: '', condition: 'new', condition_note: '', usage_months: null,
    warranty_type: 'none', warranty_months: 0, return_days: 7, media: [], hasOptions: false,
    options: [{kind: 'color', name: '', values: []}], single: {price: '', compare_at_price: '', stock_qty: 1, sku: ''}, rows: {}};
}

function fromProduct(p) {
  const L = p.source_lang;
  const e = Object.assign(emptyProduct(), {
    id: p.id, source_lang: L, name: p.translations[L]?.name || '', description: p.translations[L]?.description || '', category_id: p.category_id,
    condition: p.condition, condition_note: p.condition_note, usage_months: p.usage_months, warranty_type: p.warranty_type,
    warranty_months: p.warranty_months, return_days: p.return_days,
    media: p.media.map(m => ({key: m.key, url: m.url, kind: m.kind, poster: m.poster_url, processing: m.processing, duration: m.duration_sec})), hasOptions: p.options.length > 0,
  });
  if (p.options.length) {
    e.options = p.options.map(o => ({kind: o.kind, name: o.name[L] || '', values: o.values.map(v => ({value: v.value[L] || '', swatch_hex: v.swatch_hex || ''}))}));
    for (const v of p.variants) e.rows[rowKey(v.options)] = {on: true, price: v.price, compare_at_price: v.compare_at_price || '', stock_qty: v.stock_qty, sku: v.sku};
  } else if (p.variants[0]) {
    const v = p.variants[0];
    e.single = {price: v.price, compare_at_price: v.compare_at_price || '', stock_qty: v.stock_qty, sku: v.sku};
  }
  return e;
}

const rowKey = vals => vals.map(v => v.toLowerCase()).join('\u001f');
function combos(options) {
  return options.reduce((acc, o) => acc.flatMap(a => o.values.map(v => [...a, v.value])), [[]]);
}

async function renderProductEditor(storeId, productId) {
  const [st, cats] = await Promise.all([api(`/v1/seller/stores/${storeId}`), categories()]);
  S.store = st;
  if (!PF.E || PF.E.storeId !== storeId || PF.E.loadedId !== productId) {
    const p = productId === 'new' ? null : await api(`/v1/seller/stores/${storeId}/products/${productId}`);
    PF.E = {storeId, loadedId: productId, product: p, f: p ? fromProduct(p) : emptyProduct(), editingTr: null, step: p ? 3 : 0};
  }
  drawEditor(cats);
  if (typeof tourWizard === 'function') tourWizard(); // tour: tur wizard Tambah produk (seller-tour.js)
}

function drawEditor(cats) {
  const E = PF.E, f = E.f, p = E.product, storeId = E.storeId;
  const isNew = !p, wasSent = p && p.status !== 'draft';
  const cat = catIndex(cats)[f.category_id], needSub = !!f.category_id && cat && !cat.leaf;
  const media = kind => f.media.map((m, i) => ({...m, i})).filter(m => m.kind === kind);
  const video = media('video')[0], photos = media('photo'), defects = media('defect_photo');
  const thumb = (m, first) => `<figure class="thumb"><img src="${esc(m.url)}" alt="">${m.processing === 'processing' ? `<span class="procbadge">${t('optimizing')}</span>` : ''}
      <figcaption>${first ? '' : `<button type="button" onclick="moveMedia(${m.i})" title="${t('moveFirst')}">★</button>`}<button type="button" onclick="removeMedia(${m.i})" aria-label="${t('remove')}">✕</button></figcaption></figure>`;

  let banner = '';
  if (p) {
    const [cls, key] = productState(p, st_active());
    if (p.status === 'banned') banner = `<div class="banner bad">${t('err_product_banned')}</div>`;
    else if (p.status !== 'draft' && p.moderation_status === 'pending_review') {
      const why = (p.auto_flags?.reasons || []).map(r => `<li>${t('srs_' + r)}</li>`).join('');
      banner = `<div class="banner warn">${t('modPending')}${why ? `<div class="modwhy"><b>${t('modWhy')}</b><ul>${why}</ul></div>` : ''}</div>`;
    }
    else if (p.moderation_status === 'changes_requested') banner = `<div class="banner bad"><b>${t('modChanges')}</b> ${esc(p.moderation_note || '')}</div>`;
    else if (p.moderation_status === 'rejected') banner = `<div class="banner bad"><b>${t('modRejected')}</b> ${esc(p.moderation_note || '')}</div>`;
    else if (p.visible) banner = `<div class="banner ok">${t('modLive')}</div>`;
    banner = `<div class="row"><span class="chip ${cls}">${t(key)}</span></div>` + banner;
  }

  const optionBlocks = f.options.map((o, oi) => `<div class="card optcard">
      <div class="grid">
        <label>${t('optKind')}<select data-opt="${oi}" data-field="kind">${['color', 'size', 'material', 'capacity', 'custom'].map(k =>
          `<option value="${k}" ${o.kind === k ? 'selected' : ''}>${t('k_' + k)}</option>`).join('')}</select></label>
        <label>${t('optName')}<input data-opt="${oi}" data-field="name" value="${esc(o.name)}" maxlength="30" placeholder="${t('k_' + o.kind)}"></label>
      </div>
      <div class="row">${o.values.map((v, vi) => `<span class="vchip">${o.kind === 'color' ? `<input type="color" value="${v.swatch_hex || '#cccccc'}" data-opt="${oi}" data-val="${vi}" aria-label="${esc(v.value)}">` : ''}${esc(v.value)}
          <button type="button" onclick="removeValue(${oi},${vi})" aria-label="${t('remove')}">✕</button></span>`).join('')}
        <input class="vinput" data-addval="${oi}" placeholder="${t('valuePh')}" maxlength="40"><button type="button" onclick="addValue(${oi})">${t('addValue')}</button></div>
      ${f.options.length > 1 ? `<div><button type="button" class="danger" onclick="removeOption(${oi})">${t('removeOption')}</button></div>` : ''}
    </div>`).join('');

  const allCombos = f.hasOptions ? combos(f.options) : [];
  const ready = f.hasOptions && f.options.every(o => o.values.length);
  const variantTable = !f.hasOptions ? `<div class="grid">
      <label>${t('price')}<input id="sPrice" inputmode="decimal" value="${esc(f.single.price)}" placeholder="0.00"></label>
      <label>${t('comparePrice')} <span class="muted">(${t('optional')})</span><input id="sCompare" inputmode="decimal" value="${esc(f.single.compare_at_price)}"></label>
      <label>${t('stockQty')}<input id="sStock" type="number" min="0" value="${esc(f.single.stock_qty)}"></label>
      <label>${t('sku')} <span class="muted">(${t('optional')})</span><input id="sSku" value="${esc(f.single.sku)}" maxlength="64"></label></div>`
    : !ready ? `<p class="muted">${t('noVariants')}</p>` : `<div style="overflow-x:auto"><table class="vtable"><thead><tr><th>${t('available')}</th><th>${t('variant')}</th><th>${t('price')}</th><th>${t('comparePrice')}</th><th>${t('stockQty')}</th><th>${t('sku')}</th></tr></thead><tbody>
      <tr class="applyall"><td></td><td><b>${t('applyAll')}</b></td><td><input id="allPrice" inputmode="decimal" placeholder="0.00"></td><td><input id="allCompare" inputmode="decimal"></td><td><input id="allStock" type="number" min="0"></td><td><button type="button" onclick="applyAll()">${t('applyAll')}</button></td></tr>
      ${allCombos.map(c => { const k = rowKey(c), r = f.rows[k] || (f.rows[k] = {on: true, price: '', compare_at_price: '', stock_qty: 1, sku: ''});
        return `<tr data-row="${esc(k)}"><td><input type="checkbox" data-rf="on" ${r.on ? 'checked' : ''} aria-label="${t('available')} ${esc(c.join(' / '))}"></td><td>${esc(c.join(' / '))}</td>
          <td><input data-rf="price" inputmode="decimal" value="${esc(r.price)}" placeholder="0.00"></td><td><input data-rf="compare_at_price" inputmode="decimal" value="${esc(r.compare_at_price)}"></td>
          <td><input data-rf="stock_qty" type="number" min="0" value="${esc(r.stock_qty)}"></td><td><input data-rf="sku" value="${esc(r.sku)}" maxlength="64"></td></tr>`; }).join('')}
    </tbody></table></div>`;

  const others = p ? Object.keys(LANGS).filter(l => l !== p.source_lang) : [];
  const st = p?.translation_state || 'done';
  const trBlock = !p ? `<p class="muted">${t('translationsLater')}</p>` : others.map(l => {
    const tr = p.translations[l];
    const missing = st === 'off' ? t('trOffShort') : st === 'failed' ? t('trFailedShort') : t('translating');
    const tag = !tr ? missing : tr.seller_edited ? t('tr_seller') : tr.source === 'reviewer' ? t('tr_reviewer') : t('tr_' + tr.review_status);
    if (E.editingTr === l) return `<div class="card trcard"><b>${LANGS[l]}</b>
      <label>${t('productName')}<input id="trName" value="${esc(tr?.name || '')}" maxlength="160"></label>
      <label>${t('description')}<textarea id="trDesc">${esc(tr?.description || '')}</textarea></label>
      <div class="row"><button type="button" class="primary" onclick="saveTr('${l}')">${t('save')}</button><button type="button" onclick="PF.E.editingTr=null;drawEditor(PF.cats[P.lang])">${t('cancel')}</button></div></div>`;
    return `<div class="card trcard"><div class="spread"><b>${LANGS[l]}</b><span class="chip ${tr?.review_status === 'needs_review' || (!tr && st !== 'pending') ? 'warn' : ''}">${tag}</span></div>
      ${tr ? `<p><b>${esc(tr.name)}</b></p><p class="muted pre">${esc(tr.description)}</p>` : `<p class="muted">${missing}</p>`}
      <div class="row">${tr || st !== 'pending' ? `<button type="button" onclick="PF.E.editingTr='${l}';drawEditor(PF.cats[P.lang])">${tr ? t('editTr') : t('writeTr')}</button>` : ''}
        ${tr?.seller_edited ? `<button type="button" onclick="resetTr('${l}')">${t('useAi')}</button>` : ''}</div></div>`;
  }).join('');

  const secMedia = `<section class="card" id="secMedia">
      <div><b>${t('photos')}</b><p class="muted">${t('photosHint')}</p>${!video ? `<p class="muted small">${t('photoTips')}</p>` : ''}</div>
      <div class="thumbs" id="photoBox">${photos.map((m, n) => thumb(m, n === 0)).join('')}
        ${photos.length < 9 ? `<label class="btn drop small">+ ${t('addPhotos')}<input type="file" class="hidden" accept="image/jpeg,image/png,image/webp" multiple data-upload="photo"></label>` : ''}</div>
      <div class="vidhead"><b>${t('video')}</b> <span class="chip soft">${t('vidOptional')}</span><p class="muted">${t('videoHint')}</p></div>
      ${video ? `<div class="row"><div class="vwrap"><video class="vprev" src="${esc(video.url)}" ${video.poster ? `poster="${esc(video.poster)}"` : ''} controls preload="metadata"></video>
        ${video.processing === 'processing' ? `<span class="procbadge">${t('compressing')}</span>` : video.duration ? `<span class="procbadge ok">${video.duration}s · ${t('compressed')}</span>` : ''}</div>
          <div class="row"><label class="btn">${t('replace')}<input type="file" class="hidden" accept="video/mp4,video/quicktime,video/webm" data-upload="video"></label>
          <button type="button" class="danger" onclick="removeMedia(${video.i})">${t('remove')}</button></div></div>`
        : `<label class="btn drop" id="videoBox">▶ ${t('uploadVideo')}<input type="file" class="hidden" accept="video/mp4,video/quicktime,video/webm" data-upload="video"></label>
          <p class="vnudge">💡 ${t('noVideoNudge')}</p>
          <div class="goodvideo"><b>${t('gv_title')}</b><span>✓ ${t('gv1')}</span><span>✓ ${t('gv2')}</span><span>✓ ${t('gv3')}</span></div>`}
      ${f.condition !== 'new' ? `<div><b>${t('defects')}</b><p class="muted">${t('defectsHint')}</p></div>
        <div class="thumbs" id="defectBox">${defects.map(m => thumb(m, true)).join('')}
        ${defects.length < 6 ? `<label class="btn drop small">+ ${t('addPhotos')}<input type="file" class="hidden" accept="image/jpeg,image/png,image/webp" multiple data-upload="defect_photo"></label>` : ''}</div>` : ''}
      <p class="muted hidden" id="upState" role="status"></p>
    </section>`;
  const secBasic = `<section class="card">
      <label>${t('productName')}<input id="fName" value="${esc(f.name)}" maxlength="160" placeholder="${t('namePh')}"></label>
      <label>${t('category')}<select id="fCat"><option value="">${t('chooseCategory')}</option>${catOptions(cats, f.category_id)}</select></label>
      ${needSub ? `<div class="banner warn">${t('needSubcat')}</div>` : `<p class="muted small" style="margin-top:-6px">${t('catHint')}</p>`}
      <label>${t('description')}<textarea id="fDesc" maxlength="5000" placeholder="${t('descPh')}">${esc(f.description)}</textarea></label>
      <label>${t('writtenIn')}<select id="fLang">${Object.entries(LANGS).map(([c, n]) => `<option value="${c}" ${f.source_lang === c ? 'selected' : ''}>${n}</option>`).join('')}</select></label>
      <p class="muted">${t('translateHint')}</p>
    </section>
    <section class="card"><h2>${t('secCondition')}</h2>
      <label>${t('condition')}<select id="fCond">${['new', 'like_new', 'used_good', 'used_fair', 'refurbished'].map(c =>
        `<option value="${c}" ${f.condition === c ? 'selected' : ''} ${c !== 'new' && cat && !cat.allows_used ? 'disabled' : ''}>${t('c_' + c)}</option>`).join('')}</select></label>
      ${cat && !cat.allows_used ? `<p class="muted">${t('newOnlyCategory')}</p>` : ''}
      ${f.condition !== 'new' ? `<label>${t('conditionNote')}<textarea id="fNote" maxlength="1000" placeholder="${t('conditionNotePh')}">${esc(f.condition_note)}</textarea></label>
        <label style="max-width:240px">${t('usageMonths')} <span class="muted">(${t('optional')})</span><input id="fUsage" type="number" min="0" max="600" value="${esc(f.usage_months ?? '')}"></label>` : ''}
    </section>`;
  const secPrice = `<section class="card" id="secVariants">
      <label class="inline"><input type="checkbox" id="fHasOpt" ${f.hasOptions ? 'checked' : ''}> ${t('hasOptions')}</label>
      ${f.hasOptions ? optionBlocks + (f.options.length < 2 ? `<div><button type="button" onclick="addOption()">+ ${t('addOption')}</button></div>` : '') : ''}
      ${variantTable}
      <div class="youget" id="youGet" aria-live="polite">${ygHTML()}</div>
    </section>
    <section class="card"><h2>${t('secMore')}</h2>
      <div class="grid"><label>${t('warranty')}<select id="fWarranty">${['none', 'store', 'official'].map(w => `<option value="${w}" ${f.warranty_type === w ? 'selected' : ''}>${t('w_' + w)}</option>`).join('')}</select></label>
        ${f.warranty_type !== 'none' ? `<label>${t('warrantyMonths')}<input id="fWMonths" type="number" min="1" max="120" value="${esc(f.warranty_months || 12)}"></label>` : ''}
        <label>${t('returnDays')}<input id="fReturn" type="number" min="0" max="30" value="${esc(f.return_days)}"></label></div>
    </section>`;
  const secTr = `<section class="card"><div class="spread"><h2>${t('secTranslations')}</h2>${p && st !== 'done' ? `<span class="chip ${st === 'pending' ? 'warn' : 'bad'}">${t('trs_' + st)}</span>` : ''}</div>
      ${p ? `<p class="muted">${t('tr_source')}: ${LANGS[p.source_lang]}</p>` : ''}
      ${p && st === 'off' ? `<div class="banner warn">${t('trOff')}</div>` : ''}
      ${p && (st === 'failed' || st === 'retrying') ? `<div class="banner bad"><b>${st === 'failed' ? t('trFailed') : t('trRetrying')}</b>
        ${p.translation_error ? `<br><span class="muted" style="color:inherit">${esc(p.translation_error)}</span>` : ''}
        <div style="margin-top:8px"><button type="button" onclick="retryTr()">${t('trRetry')}</button></div></div>` : ''}
      <div class="grid">${trBlock}</div></section>`;
  // Langkah 4: pratinjau + daftar periksa (tiap baris bisa "Ubah" → kembali ke langkahnya).
  const cover = photos[0]?.url || video?.poster || '';
  const prices = f.hasOptions ? Object.values(f.rows).filter(r => r.on).map(r => Number(r.price)).filter(n => n > 0) : [Number(f.single.price)].filter(n => n > 0);
  const stockAll = f.hasOptions ? Object.values(f.rows).filter(r => r.on).reduce((a, r) => a + (Number(r.stock_qty) || 0), 0) : Number(f.single.stock_qty) || 0;
  const priceTxt = !prices.length ? '—' : Math.min(...prices) === Math.max(...prices) ? money(prices[0]) : `${money(Math.min(...prices))} – ${money(Math.max(...prices))}`;
  const rv = (s, key) => { const ok = stepOk(s, true); return `<div class="rv ${ok ? '' : 'bad'}"><i aria-hidden="true">${ok ? '✓' : '!'}</i><span>${t(key)} · <span class="muted">${t(ok ? 'rv_ok' : 'rv_bad')}</span></span>
      <button type="button" onclick="wzGo(${s})">${t('wzEdit')}</button></div>`; };
  const secCheck = `<section class="card"><div class="preview">${cover ? `<img class="pv" src="${esc(cover)}" alt="">` : '<span class="pv" aria-hidden="true"></span>'}
      <div><b style="font-size:17px">${esc(f.name || '—')}</b><p class="muted">${esc(cat?.name || '—')} · ${t('c_' + f.condition)}</p>
        <p style="font-size:20px;font-weight:800;color:var(--pink-text)">${priceTxt}</p><p class="muted">${t('stock')} ${stockAll}</p></div></div></section>
    <section class="card reviewlist">${rv(0, 'rv_media')}${rv(1, 'rv_basic')}${rv(2, 'rv_price')}</section>
    ${wasSent ? `<p class="muted">${t('reviewNote')}</p>` : ''}
    ${p ? `<div class="row">
      ${p && p.status !== 'draft' && p.status !== 'banned' ? `<button type="button" onclick="setPStatus('${p.status === 'archived' ? 'active' : 'archived'}')">${p.status === 'archived' ? t('unarchive') : t('archive')}</button>` : ''}
      ${p ? `<button type="button" class="danger" onclick="deleteProduct()">${t('del')}</button>` : ''}</div>` : ''}
    ${p ? secTr : ''}`;
  const keys = ['media', 'basic', 'price', 'check'], step = E.step || 0;
  const bar = keys.map((k, i) => `<button type="button" class="${i < step ? 'done' : i === step ? 'on' : ''}" onclick="wzGo(${i})" ${i === step ? 'aria-current="step"' : ''}><i></i>${t('wz_' + k)}</button>`).join('');
  const body = [secMedia, secBasic, secPrice, secCheck][step];
  const foot = step < 3
    ? `${step > 0 ? `<button type="button" onclick="wzGo(${step - 1})">${t('wzBack')}</button>` : ''}${!wasSent ? `<button type="button" id="btnDraft">${t('saveDraft')}</button>` : ''}
       <span class="grow"></span><button type="button" class="primary" onclick="wzNext()">${t('wzNext', {s: t('wz_' + keys[step + 1])})}</button>`
    : !wasSent && S.store?.status === 'onboarding' && S.store?.role === 'owner'
    // Toko baru (belum diverifikasi): langkah terakhir = simpan draf lalu langsung ke verifikasi toko.
    ? `<button type="button" onclick="wzGo(2)">${t('wzBack')}</button><span class="grow"></span>
       <button type="button" class="primary" id="btnVerify">${t('draftVerify')} →</button>
       <span class="muted" style="flex-basis:100%">${t('draftVerifyHint')}</span>`
    : `<button type="button" onclick="wzGo(2)">${t('wzBack')}</button>${!wasSent ? `<button type="button" id="btnDraft">${t('saveDraft')}</button>` : ''}<span class="grow"></span>
       <button type="button" class="primary" id="btnSubmit" ${p?.status === 'banned' || (!wasSent && !st_active()) ? 'disabled' : ''}>${wasSent ? t('saveResubmit') : t('submitReview')}</button>
       ${!wasSent && !st_active() ? `<span class="muted" style="flex-basis:100%">${t(S.store?.status === 'pending_review' ? 'storeReviewing' : 'submitLocked')}</span>` : ''}`;

  $('#app').innerHTML = focusShell(S.store, `
    <a class="backlink" href="#/store/${storeId}/products">${t('backProducts')}</a>
    <div class="wiz"><p class="muted" style="font-weight:700">${t('wzStepOf', {a: step + 1})}</p>
      <h1>${isNew ? t('addProduct') : t('editProduct')} · ${t('wz_' + keys[step])}</h1>
      <div class="wizbar">${bar}</div><p class="wizhint">${t('wzh_' + keys[step])}</p></div>
    ${banner}
    ${st_active() ? '' : `<div class="banner info">${t('storeNotActiveHint')}</div>`}
    ${body}
    <div class="wizfoot" data-errors-anchor>${foot}</div>`);
  bindEditor();
  clearTimeout(PF.poll);
  const mediaBusy = p && p.media.some(m => m.processing === 'processing');
  if (p && (st === 'pending' || st === 'retrying' || mediaBusy)) PF.poll = setTimeout(pollTranslations, st === 'pending' || mediaBusy ? 4000 : 15000);
}

const st_active = () => S.store?.status === 'active';

// Simpan isi form ke PF.E.f sebelum menggambar ulang.
function readForm() {
  const f = PF.E.f, v = id => $('#' + id)?.value;
  if ($('#fName')) Object.assign(f, {name: v('fName'), description: v('fDesc'), source_lang: v('fLang'), category_id: v('fCat'), condition: v('fCond')});
  if ($('#fNote')) f.condition_note = v('fNote');
  if ($('#fUsage')) f.usage_months = v('fUsage') === '' ? null : Number(v('fUsage'));
  if ($('#fWarranty')) f.warranty_type = v('fWarranty');
  if ($('#fWMonths')) f.warranty_months = Number(v('fWMonths'));
  if ($('#fReturn')) f.return_days = Number(v('fReturn'));
  if ($('#sPrice')) f.single = {price: v('sPrice'), compare_at_price: v('sCompare'), stock_qty: Number(v('sStock')), sku: v('sSku')};
  document.querySelectorAll('[data-opt][data-field]').forEach(el => { f.options[el.dataset.opt][el.dataset.field] = el.value; });
  document.querySelectorAll('[data-opt][data-val]').forEach(el => { const v = f.options[el.dataset.opt].values[el.dataset.val]; if (v.swatch_hex || el.value !== '#cccccc') v.swatch_hex = el.value; });
  document.querySelectorAll('tr[data-row]').forEach(tr => {
    const r = f.rows[tr.dataset.row];
    tr.querySelectorAll('[data-rf]').forEach(el => { r[el.dataset.rf] = el.type === 'checkbox' ? el.checked : el.dataset.rf === 'stock_qty' ? Number(el.value) : el.value; });
  });
}
const redraw = () => { readForm(); drawEditor(PF.cats[P.lang]); };

function bindEditor() {
  ['fCond', 'fCat', 'fWarranty'].forEach(id => $('#' + id) && ($('#' + id).onchange = redraw));
  if ($('#fHasOpt')) $('#fHasOpt').onchange = e => { readForm(); PF.E.f.hasOptions = e.target.checked; drawEditor(PF.cats[P.lang]); };
  document.querySelectorAll('[data-opt][data-field=kind]').forEach(el => el.onchange = redraw);
  document.querySelectorAll('[data-addval]').forEach(el => el.onkeydown = e => { if (e.key === 'Enter') { e.preventDefault(); addValue(Number(el.dataset.addval)); } });
  document.querySelectorAll('[data-upload]').forEach(el => el.onchange = e => uploadMedia(el.dataset.upload, [...e.target.files]));
  if ($('#btnDraft')) $('#btnDraft').onclick = () => saveProduct(false);
  if ($('#btnSubmit')) $('#btnSubmit').onclick = () => saveProduct(true);
  if ($('#btnVerify')) $('#btnVerify').onclick = draftThenVerify;
  document.querySelectorAll('#sPrice, [data-rf=price], tr[data-row] [data-rf=on]').forEach(el => el.addEventListener('input', () => { readForm(); const y = $('#youGet'); if (y) y.innerHTML = ygHTML(); }));
}

function addValue(oi) {
  readForm();
  const inp = document.querySelector(`[data-addval="${oi}"]`), o = PF.E.f.options[oi];
  for (const raw of inp.value.split(',')) {
    const val = raw.trim();
    if (val && !o.values.some(v => v.value.toLowerCase() === val.toLowerCase())) o.values.push({value: val, swatch_hex: ''});
  }
  drawEditor(PF.cats[P.lang]);
  document.querySelector(`[data-addval="${oi}"]`)?.focus();
}
function removeValue(oi, vi) { readForm(); PF.E.f.options[oi].values.splice(vi, 1); drawEditor(PF.cats[P.lang]); }
function addOption() { readForm(); PF.E.f.options.push({kind: PF.E.f.options[0]?.kind === 'size' ? 'color' : 'size', name: '', values: []}); drawEditor(PF.cats[P.lang]); }
function removeOption(oi) { readForm(); PF.E.f.options.splice(oi, 1); drawEditor(PF.cats[P.lang]); }
function applyAll() {
  readForm();
  const price = $('#allPrice').value, cmp = $('#allCompare').value, stock = $('#allStock').value;
  for (const r of Object.values(PF.E.f.rows)) {
    if (price) r.price = price;
    if (cmp) r.compare_at_price = cmp;
    if (stock !== '') r.stock_qty = Number(stock);
  }
  drawEditor(PF.cats[P.lang]);
}
function removeMedia(i) { readForm(); PF.E.f.media.splice(i, 1); drawEditor(PF.cats[P.lang]); }
function moveMedia(i) {
  readForm();
  const m = PF.E.f.media.splice(i, 1)[0], first = PF.E.f.media.findIndex(x => x.kind === m.kind);
  PF.E.f.media.splice(first < 0 ? PF.E.f.media.length : first, 0, m);
  drawEditor(PF.cats[P.lang]);
}

function videoDuration(file) {
  return new Promise(res => {
    const v = document.createElement('video');
    v.preload = 'metadata';
    v.onloadedmetadata = () => { URL.revokeObjectURL(v.src); res(v.duration); };
    v.onerror = () => res(0); // format tidak bisa dibaca browser (mis. HEVC) → biar server yang menilai
    v.src = URL.createObjectURL(file);
  });
}

// Upload dengan progress (video bisa besar). Token diperbarui sekali kalau kedaluwarsa.
function xhrUpload(file, onProgress, retried) {
  return new Promise((resolve, reject) => {
    const fd = new FormData(); fd.append('kind', 'product_media'); fd.append('file', file);
    const x = new XMLHttpRequest();
    x.open('POST', '/v1/seller/uploads');
    x.setRequestHeader('Authorization', 'Bearer ' + P.token);
    x.setRequestHeader('X-Lang', P.lang);
    x.upload.onprogress = e => e.lengthComputable && onProgress(Math.round(e.loaded / e.total * 100));
    x.onload = async () => {
      let data = {}; try { data = JSON.parse(x.responseText); } catch {}
      if (x.status === 401 && !retried && await refreshSession()) return resolve(xhrUpload(file, onProgress, true));
      if (x.status >= 200 && x.status < 300) resolve(data);
      else reject(Object.assign(new Error(data.error?.message || x.statusText), {code: data.error?.code}));
    };
    x.onerror = () => reject(new Error('network'));
    x.send(fd);
  });
}

async function uploadMedia(kind, files) {
  readForm();
  const state = $('#upState');
  for (const file of files) {
    if (kind === 'video') {
      if (file.size > 100 * 1024 * 1024) { toast(t('fileTooBig'), 'warning'); continue; }
      if (await videoDuration(file) > 181) { toast(t('videoTooLong'), 'warning'); continue; }
    } else if (file.size > 8 * 1024 * 1024) { toast(t('fileTooBig'), 'warning'); continue; }
    if (kind === 'photo' && PF.E.f.media.filter(m => m.kind === 'photo').length >= 9) break;
    state.classList.remove('hidden');
    try {
      const up = await xhrUpload(file, pct => { state.textContent = t('uploading', {p: pct}); });
      if (kind === 'video') PF.E.f.media = PF.E.f.media.filter(m => m.kind !== 'video');
      const item = {key: up.key, url: up.url, kind};
      if (kind === 'video') PF.E.f.media.unshift(item); else PF.E.f.media.push(item);
    } catch (e) { toast(errMsg(e), 'error'); }
  }
  drawEditor(PF.cats[P.lang]);
}

function buildBody(submit) {
  const f = PF.E.f;
  const body = {category_id: f.category_id, source_lang: f.source_lang, name: f.name, description: f.description, condition: f.condition,
    condition_note: f.condition === 'new' ? '' : f.condition_note, usage_months: f.condition === 'new' ? null : f.usage_months,
    warranty_type: f.warranty_type, warranty_months: f.warranty_type === 'none' ? 0 : f.warranty_months, return_days: f.return_days,
    media: f.media.filter(m => f.condition !== 'new' || m.kind !== 'defect_photo').map(m => ({key: m.key, kind: m.kind})), submit,
    options: [], variants: []};
  PF.E.sentRows = [];
  if (f.hasOptions) {
    body.options = f.options.map(o => ({kind: o.kind, name: o.name.trim() || t('k_' + o.kind), values: o.values.map(v => ({value: v.value, swatch_hex: o.kind === 'color' ? (v.swatch_hex || '') : ''}))}));
    for (const c of combos(f.options)) {
      const r = f.rows[rowKey(c)];
      if (r && r.on) PF.E.sentRows.push(rowKey(c));
      if (r && r.on) body.variants.push({options: c, price: String(r.price).trim(), compare_at_price: String(r.compare_at_price || '').trim(), stock_qty: Number(r.stock_qty) || 0, sku: (r.sku || '').trim()});
    }
  } else {
    body.variants = [{options: [], price: String(f.single.price).trim(), compare_at_price: String(f.single.compare_at_price || '').trim(), stock_qty: Number(f.single.stock_qty) || 0, sku: (f.single.sku || '').trim()}];
  }
  return body;
}

// Wizard: pindah langkah (isian disimpan dulu), "Lanjut" memeriksa langkah ini dulu.
function wzGo(i) { readForm(); clearFieldErrors(); PF.E.step = i; drawEditor(PF.cats[P.lang]); window.scrollTo(0, 0); }
function wzNext() {
  readForm();
  const s = PF.E.step || 0;
  if (!stepOk(s, true)) { checkRequired(stepRules(s, true)); return; }
  wzGo(s + 1);
}

// Tarif komisi yang berlaku untuk produk ini (toko > kategori > global) dan teks "Anda terima".
function prodRate() {
  const c = S.store?.commission || {};
  return parseFloat(c.store || (c.categories || {})[PF.E.f.category_id] || c.global || 8);
}
function ygHTML() {
  const f = PF.E.f, r = prodRate();
  const prices = f.hasOptions ? Object.values(f.rows).filter(x => x.on).map(x => Number(x.price)).filter(n => n > 0) : [Number(f.single.price)].filter(n => n > 0);
  if (!prices.length) return `<span>${t('yg_title')}</span><small>${t('yg_none')}</small>`;
  const p = Math.min(...prices), fee = Math.round(p * r) / 100;
  return `<span>${prices.length > 1 && Math.max(...prices) !== p ? t('yg_from') : t('yg_title')}</span><b>${money(p - fee)}</b>
    <small>${t('yg_line', {p: p.toFixed(2), r, f: fee.toFixed(2)})}</small>`;
}

// Video opsional: tanpa video wajib ≥3 foto (sama dengan minPhotosNoVideo di backend).
const MIN_PHOTOS_NO_VIDEO = 3;
function photoOk(kinds) {
  const n = kinds.filter(k => k === 'photo').length;
  return n > 0 && (kinds.includes('video') || n >= MIN_PHOTOS_NO_VIDEO);
}

// Cek berbasis data (tanpa DOM) — dipakai untuk daftar periksa & untuk tahu langkah mana yang kurang.
function stepOk(step, submit) {
  const f = PF.E.f, kinds = f.media.map(m => m.kind), pos = v => Number(String(v ?? '').trim()) > 0;
  const cmpOk = (c, p) => String(c ?? '').trim() === '' || Number(c) > Number(p);
  if (step === 0) return !submit || (photoOk(kinds) && (!['used_good', 'used_fair'].includes(f.condition) || kinds.includes('defect_photo')));
  if (step === 1) return f.name.trim().length >= 3 && !!catIndex(PF.cats[P.lang])[f.category_id]?.leaf && (f.condition === 'new' || String(f.condition_note || '').trim() !== '');
  if (step === 2) {
    if (!f.hasOptions) return pos(f.single.price) && cmpOk(f.single.compare_at_price, f.single.price);
    if (!f.options.every(o => o.values.length)) return false;
    const on = combos(f.options).map(c => f.rows[rowKey(c)]).filter(r => r && r.on);
    return on.length > 0 && on.every(r => pos(r.price) && cmpOk(r.compare_at_price, r.price));
  }
  return true;
}

// Aturan DOM untuk langkah yang sedang tampil (menandai isian yang salah).
function stepRules(step, submit) {
  const f = PF.E.f, kinds = f.media.map(m => m.kind);
  const price = el => Number(String(el.value).trim()) > 0;
  const compareOk = (cmpEl, priceEl) => { const c = String(cmpEl.value).trim(); return c === '' || Number(c) > Number(priceEl.value); };
  const rules = [];
  if (step === 0 && submit) {
    rules.push(['#photoBox', t('needPhoto'), () => kinds.includes('photo'), t('photos')]);
    rules.push(['#photoBox', t('needMorePhotos'), () => photoOk(kinds), t('photos')]);
    if (['used_good', 'used_fair'].includes(f.condition)) rules.push(['#defectBox', t('needDefect'), () => kinds.includes('defect_photo'), t('defects')]);
  }
  if (step === 1) {
    rules.push(['#fName', t('minChars', {n: 3}), el => el.value.trim().length >= 3], ['#fCat', t('pickOne')]);
    if (f.condition !== 'new') rules.push(['#fNote']);
  }
  if (step === 2) {
    if (!f.hasOptions) {
      rules.push(['#sPrice', t('positive'), price]);
      rules.push(['#sCompare', t('compareHigher'), el => compareOk(el, $('#sPrice'))]);
    } else {
      f.options.forEach((o, i) => rules.push([`[data-addval="${i}"]`, t('needValue'), () => o.values.length > 0, `${t('option')} ${i + 1} (${t('k_' + o.kind)})`]));
      rules.push(['#secVariants', t('needVariantOn'), () => !document.querySelector('tr[data-row]') || !!document.querySelector('tr[data-row] [data-rf=on]:checked'), t('secVariants')]);
      document.querySelectorAll('tr[data-row]').forEach(tr => {
        if (!tr.querySelector('[data-rf=on]').checked) return;
        rules.push([tr.querySelector('[data-rf=price]'), t('positive'), price]);
        rules.push([tr.querySelector('[data-rf=compare_at_price]'), t('compareHigher'), el => compareOk(el, tr.querySelector('[data-rf=price]'))]);
      });
    }
  }
  return rules;
}

// Sebelum menyimpan: cari langkah pertama yang kurang, buka langkah itu, tandai isiannya.
function validateProductForm(submit) {
  for (const s of [0, 1, 2]) {
    if (stepOk(s, submit)) continue;
    if (PF.E.step !== s) { PF.E.step = s; drawEditor(PF.cats[P.lang]); }
    checkRequired(stepRules(s, submit));
    return false;
  }
  return true;
}

// Error dari server → langkah tempat isian itu berada.
function fieldStep(field) {
  if (!field) return null;
  if (/^media/.test(field)) return 0;
  if (/^(name|description|category_id|source_lang|condition|condition_note|usage_months)$/.test(field)) return 1;
  if (/^(options|variants|warranty_type|warranty_months|return_days)/.test(field)) return 2;
  return null;
}

function productErrorField(field) {
  const simple = {name: '#fName', description: '#fDesc', category_id: '#fCat', source_lang: '#fLang', condition: '#fCond',
    condition_note: '#fNote', usage_months: '#fUsage', warranty_type: '#fWarranty', warranty_months: '#fWMonths', return_days: '#fReturn',
    options: '#secVariants', variants: '#secVariants', media: '#secMedia'};
  if (simple[field]) return simple[field];
  let m = field.match(/^options\[(\d+)\]\.(\w+)/);
  if (m) return m[2] === 'values' ? `[data-addval="${m[1]}"]` : `[data-opt="${m[1]}"][data-field="${m[2]}"]`;
  m = field.match(/^options\[(\d+)\]/);
  if (m) return `[data-addval="${m[1]}"]`;
  m = field.match(/^variants\[(\d+)\]\.(\w+)/);
  if (m) {
    const col = m[2] === 'options' ? 'on' : m[2];
    if (!PF.E.f.hasOptions) return ({price: '#sPrice', compare_at_price: '#sCompare', stock_qty: '#sStock', sku: '#sSku'})[col];
    const key = PF.E.sentRows?.[Number(m[1])];
    return key ? document.querySelector(`tr[data-row="${CSS.escape(key)}"] [data-rf="${col}"]`) : '#secVariants';
  }
  if (/^media\[(\d+)\]/.test(field)) return '#secMedia';
  return null;
}

async function saveProduct(submit) {
  readForm();
  if (!validateProductForm(submit)) return;
  const E = PF.E, url = `/v1/seller/stores/${E.storeId}/products`;
  // Harga sangat mahal biasanya salah ketik (mis. 15000 maksudnya 15.00) — tanya dulu, bukan ditolak.
  const top = Math.max(0, ...buildBody(submit).variants.map(v => Number(v.price) || 0));
  if (submit && top >= 5000 && !confirm(t('priceHuge', {p: money(top)}))) return;
  try {
    const p = E.product ? await api(`${url}/${E.product.id}`, {method: 'PUT', body: buildBody(submit)})
      : await api(url, {method: 'POST', body: buildBody(submit)});
    toast(!submit ? t('savedDraft') : p.moderation_status === 'approved' ? t('nowLive') : t('sentReview'));
    PF.E = {storeId: E.storeId, loadedId: p.id, product: p, f: fromProduct(p), editingTr: null, step: 3};
    if (location.hash !== `#/store/${E.storeId}/product/${p.id}`) history.replaceState(null, '', `#/store/${E.storeId}/product/${p.id}`);
    drawEditor(PF.cats[P.lang]);
    return p;
  } catch (e) {
    const fs = fieldStep(e.field);
    if (fs !== null && fs !== PF.E.step) { PF.E.step = fs; drawEditor(PF.cats[P.lang]); }
    showError(e, productErrorField);
  }
}

// Toko baru: simpan produk sebagai draf, lalu langsung ke verifikasi toko — kalau semua langkah sudah lengkap
// toko langsung dikirim untuk ditinjau, kalau belum seller dibawa ke langkah yang masih kurang.
async function draftThenVerify() {
  const btn = $('#btnVerify');
  if (btn) btn.disabled = true;
  const p = await saveProduct(false);
  if (!p) { if ($('#btnVerify')) $('#btnVerify').disabled = false; return; }
  const id = PF.E.storeId;
  try {
    const st = await api(`/v1/seller/stores/${id}`);
    S.store = st;
    const next = setupSteps(st).find(x => x.n && !x.done);
    if (!next || next.k === 'submit') {
      if (st.checklist?.can_submit) {
        await api(`/v1/seller/stores/${id}/submit`, {method: 'POST'});
        toast(t('sw_sentT'));
      }
      S.tab = 'home';
    } else {
      toast(t('draftVerifyFix'), 'info');
      S.tab = next.k;
    }
    try { sessionStorage.setItem('phsar-tab-' + id, S.tab); } catch {}
    S.storeId = id;
    location.hash = '#/store/' + id;
    window.scrollTo(0, 0);
  } catch (e) { toast(errMsg(e), 'error'); if ($('#btnVerify')) $('#btnVerify').disabled = false; }
}

async function refreshProduct() {
  const E = PF.E;
  readForm();
  const p = await api(`/v1/seller/stores/${E.storeId}/products/${E.product.id}`);
  E.product = p;
  // Hasil kompres/sampul yang baru selesai → perbarui pratinjau tanpa mengubah isian form.
  for (const m of E.f.media) {
    const fresh = p.media.find(x => x.key === m.key);
    if (fresh) Object.assign(m, {url: fresh.url, poster: fresh.poster_url, processing: fresh.processing, duration: fresh.duration_sec});
  }
  drawEditor(PF.cats[P.lang]);
}
async function pollTranslations() {
  if (!PF.E?.product || !location.hash.includes('/product/' + PF.E.product.id) || PF.E.editingTr) return;
  try { await refreshProduct(); } catch {}
}
async function saveTr(lang) {
  const E = PF.E;
  if (!checkRequired([['#trName']])) return;
  try {
    E.product = await api(`/v1/seller/stores/${E.storeId}/products/${E.product.id}/translations/${lang}`, {method: 'PUT', body: {name: $('#trName').value, description: $('#trDesc').value}});
    E.editingTr = null; toast(t('saved')); drawEditor(PF.cats[P.lang]);
  } catch (e) { showError(e, {name: '#trName', description: '#trDesc'}); }
}
async function retryTr() {
  const E = PF.E;
  try { readForm(); E.product = await api(`/v1/seller/stores/${E.storeId}/products/${E.product.id}/translations/retry`, {method: 'POST'}); drawEditor(PF.cats[P.lang]); }
  catch (e) { toast(errMsg(e), 'error'); }
}
async function resetTr(lang) {
  const E = PF.E;
  try { readForm(); E.product = await api(`/v1/seller/stores/${E.storeId}/products/${E.product.id}/translations/${lang}`, {method: 'DELETE'}); drawEditor(PF.cats[P.lang]); }
  catch (e) { toast(errMsg(e), 'error'); }
}
async function setPStatus(status) {
  const E = PF.E;
  try { readForm(); E.product = await api(`/v1/seller/stores/${E.storeId}/products/${E.product.id}/status`, {method: 'POST', body: {status}}); drawEditor(PF.cats[P.lang]); }
  catch (e) { toast(errMsg(e), 'error'); }
}
async function deleteProduct() {
  const E = PF.E;
  if (!confirm(t('confirmDelete'))) return;
  try { await api(`/v1/seller/stores/${E.storeId}/products/${E.product.id}`, {method: 'DELETE'}); PF.E = null; location.hash = '#/store/' + E.storeId; }
  catch (e) { toast(errMsg(e), 'error'); }
}
