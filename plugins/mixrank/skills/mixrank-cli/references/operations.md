# MixRank operation reference

Generated from catalog version 2026-09-09.

102 operations visible in the authorized documentation, including five Elasticsearch searches. Licensing still determines which calls may succeed.

Use `mixrank api <command> --help` for flags. All parameters are strings on the wire; JSON bodies preserve numeric precision. `--query KEY=VALUE` passes additional documented query fields.

## delete-email-validate-bulk-job-by-job-id

`DELETE /email/validate/bulk-job/{job_id}` — Email » Delete Job Output

Encoding: query. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `job_id` | path | true |  |

Response fields: `id` (str), `created_at` (datetime), `lines_total` (int), `lines_invalid` (int), `lines_processed` (int), `lines_refreshed` (int), `strategy` (str), `maxage` (datetime), `completed_at` (datetime), `download_url` (str), `private` (bool).

Sources: [provider contract](https://mixrank.com/api/documentation#/email/validate/bulk-job/{job_id}). Verified 2026-09-09 signed-in documentation.

## get-appstore-apps

`GET /appstore/apps` — iOS Apps » Directory

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=rank |
| `sort_order` | query | false | default=desc |
| `company.id` | query | false |  |
| `search` | query | false |  |

Response fields: `id` (str), `bundle_id` (str), `title` (str), `rank` (int), `score` (float), `developer_id` (int), `developer` (str), `company_url` (str), `support_url` (str), `privacy_policy_url` (str), `seller` (str), `copyright` (str), `downloaded` (bool), `review_count` (int), `rating` (float), `rating_count` (int), `upload_date` (date), `release_date` (date), `small_artwork_url` (str), `large_artwork_url` (str), `xlarge_artwork_url` (str), `genre` (str), `genres` (str), `configs_count` (int), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `version` (object), `version.version_number` (str), `version.version_string` (str), `content_rating` (str), `content_rating_details` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/apps). Verified 2026-09-09 signed-in documentation.

## get-appstore-apps-by-app-id

`GET /appstore/apps/{app_id}` — iOS Apps » Overview

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `app_id` | path | true |  |

Response fields: `id` (str), `bundle_id` (str), `title` (str), `rank` (int), `score` (float), `developer_id` (int), `developer` (str), `company_url` (str), `support_url` (str), `privacy_policy_url` (str), `seller` (str), `copyright` (str), `downloaded` (bool), `review_count` (int), `rating` (float), `rating_count` (int), `upload_date` (date), `release_date` (date), `small_artwork_url` (str), `large_artwork_url` (str), `xlarge_artwork_url` (str), `genre` (str), `genres` (str), `configs_count` (int), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `version` (object), `version.version_number` (str), `version.version_string` (str), `content_rating` (str), `content_rating_details` (str), `developer_link` (str), `ratings` (object), `ratings.1` (int), `ratings.2` (int), `ratings.3` (int), `ratings.4` (int), `ratings.5` (int), `summary` (str), `installation_size` (int), `price` (str), `compatibility` (object), `compatibility.min_ios_version` (str), `compatibility.supported_devices` (list[str]), `screenshot_urls` (list[str]), `ipad_screenshot_urls` (list[str]), `privacy_label_ids` (list[int]).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/apps/{app_id}). Verified 2026-09-09 signed-in documentation.

## get-appstore-apps-by-app-id-configs

`GET /appstore/apps/{app_id}/configs` — iOS Apps » Configs

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `app_id` | path | true |  |

Response fields: `version` (str), `plist` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/apps/{app_id}/configs). Verified 2026-09-09 signed-in documentation.

## get-appstore-apps-by-app-id-downloads

`GET /appstore/apps/{app_id}/downloads` — iOS Apps » Downloads

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `app_id` | path | true |  |
| `interval` | query | false | default=1 month |
| `country` | query | false |  |
| `offset` | query | false |  |

Response fields: `country` (str), `data` (list[list]).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/apps/{app_id}/downloads). Verified 2026-09-09 signed-in documentation.

## get-appstore-apps-by-app-id-iaps

`GET /appstore/apps/{app_id}/iaps` — iOS Apps » In-App Purchases

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `app_id` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=id |
| `sort_order` | query | false | default=desc |

Response fields: `id` (int), `name` (str), `price_str` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/apps/{app_id}/iaps). Verified 2026-09-09 signed-in documentation.

## get-appstore-apps-by-app-id-rankings

`GET /appstore/apps/{app_id}/rankings` — iOS Apps » Rankings Summary

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `app_id` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |

Response fields: `id` (int), `container` (str), `genre` (str), `country` (str), `current_rank` (int), `highest_rank` (int), `first_ranked` (date), `last_ranked` (date).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/apps/{app_id}/rankings). Verified 2026-09-09 signed-in documentation.

## get-appstore-apps-by-app-id-rankings-history

`GET /appstore/apps/{app_id}/rankings-history` — iOS Apps » Rankings History

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `app_id` | path | true |  |
| `container` | query | true |  |
| `genre` | query | true |  |
| `country` | query | true |  |
| `before` | query | false |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |

Response fields: `ranking` (int), `date` (datetime).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/apps/{app_id}/rankings-history). Verified 2026-09-09 signed-in documentation.

## get-appstore-apps-by-app-id-revenue

`GET /appstore/apps/{app_id}/revenue` — iOS Apps » Revenue

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `app_id` | path | true |  |
| `interval` | query | false | default=1 month |
| `country` | query | false |  |
| `offset` | query | false |  |

Response fields: `country` (str), `data` (list[list]).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/apps/{app_id}/revenue). Verified 2026-09-09 signed-in documentation.

## get-appstore-apps-by-app-id-reviews

`GET /appstore/apps/{app_id}/reviews` — iOS Apps » Reviews

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `app_id` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=date |
| `sort_order` | query | false | default=desc |

Response fields: `id` (int), `user_id` (int), `username` (str), `rating` (int), `version` (str), `date` (date), `vote_yes` (int), `vote_no` (int), `title` (str), `body` (str), `country_id` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/apps/{app_id}/reviews). Verified 2026-09-09 signed-in documentation.

## get-appstore-apps-by-app-id-sdks

`GET /appstore/apps/{app_id}/sdks` — iOS Apps » SDKs

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `app_id` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=500 |
| `sort_field` | query | false | default=version_code_min |
| `sort_order` | query | false | default=desc |

Response fields: `name` (str), `slug` (str), `upload_date_min` (date), `upload_date_max` (date), `version_code_min` (int), `version_code_max` (int), `currently_installed` (bool), `sdk_version` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/apps/{app_id}/sdks). Verified 2026-09-09 signed-in documentation.

## get-appstore-apps-by-app-id-versions

`GET /appstore/apps/{app_id}/versions` — iOS Apps » Versions

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `app_id` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=id |
| `sort_order` | query | false | default=desc |

Response fields: `id` (int), `asset_size` (int), `config_count` (int), `class_count` (int), `resource_count` (int), `version_display` (str), `release_date` (date), `release_notes` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/apps/{app_id}/versions). Verified 2026-09-09 signed-in documentation.

## get-appstore-developers-by-developer-id-apps

`GET /appstore/developers/{developer_id}/apps` — iOS Developers » Apps

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `developer_id` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=rating_count |
| `sort_order` | query | false | default=desc |

Response fields: `id` (str), `bundle_id` (str), `title` (str), `rank` (int), `score` (float), `developer_id` (int), `developer` (str), `company_url` (str), `support_url` (str), `privacy_policy_url` (str), `seller` (str), `copyright` (str), `downloaded` (bool), `review_count` (int), `rating` (float), `rating_count` (int), `upload_date` (date), `release_date` (date), `small_artwork_url` (str), `large_artwork_url` (str), `xlarge_artwork_url` (str), `genre` (str), `genres` (str), `configs_count` (int), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `version` (object), `version.version_number` (str), `version.version_string` (str), `content_rating` (str), `content_rating_details` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/developers/{developer_id}/apps). Verified 2026-09-09 signed-in documentation.

## get-appstore-ipa-by-ipa-id

`GET /appstore/ipa/{ipa_id}` — iOS Apps » IPA

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `ipa_id` | path | true |  |

Response fields: `ipa_id` (int), `error` (str), `download_link` (str), `app_id` (int), `asset_size` (int), `bundle_short_version_string` (str), `bundle_version` (str), `configs_count` (int), `md5` (str), `resource_count` (int), `version_display` (str), `downloaded_at` (datetime), `apple_platform_name` (str), `release_date` (datetime), `release_notes` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/ipa/{ipa_id}). Verified 2026-09-09 signed-in documentation.

## get-appstore-privacy-datatype

`GET /appstore/privacy/datatype` — iOS App Privacy » Privacy Data Types

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



Response fields: `id` (int), `category_label` (str), `category_identifier` (str), `label` (str), `description` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/privacy/datatype). Verified 2026-09-09 signed-in documentation.

## get-appstore-privacy-labels

`GET /appstore/privacy/labels` — iOS App Privacy » Privacy labels

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



Response fields: `id` (int), `datatype_id` (int), `linkage_id` (int), `purpose_id` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/privacy/labels). Verified 2026-09-09 signed-in documentation.

## get-appstore-privacy-linkage

`GET /appstore/privacy/linkage` — iOS App Privacy » Privacy Linkage Types

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



Response fields: `id` (int), `label` (str), `identifier` (str), `description` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/privacy/linkage). Verified 2026-09-09 signed-in documentation.

## get-appstore-privacy-purpose

`GET /appstore/privacy/purpose` — iOS App Privacy » Privacy Data Purposes

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



Response fields: `id` (int), `label` (str), `identifier` (str), `description` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/privacy/purpose). Verified 2026-09-09 signed-in documentation.

## get-appstore-rankings

`GET /appstore/rankings` — iOS Rankings » Lists

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



Response fields: `containers` (list[str]), `genres` (list[str]), `countries` (list[str]).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/rankings). Verified 2026-09-09 signed-in documentation.

## get-appstore-rankings-by-container-by-genre-by-country

`GET /appstore/rankings/{container}/{genre}/{country}` — iOS Rankings » Rankings

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=ranking |
| `sort_order` | query | false | default=asc |
| `container` | path | true |  |
| `genre` | path | true |  |
| `country` | path | true |  |
| `before` | query | false |  |

Response fields: `id` (str), `bundle_id` (str), `title` (str), `rank` (int), `score` (float), `developer_id` (int), `developer` (str), `company_url` (str), `support_url` (str), `privacy_policy_url` (str), `seller` (str), `copyright` (str), `downloaded` (bool), `review_count` (int), `rating` (float), `rating_count` (int), `upload_date` (date), `release_date` (date), `small_artwork_url` (str), `large_artwork_url` (str), `xlarge_artwork_url` (str), `genre` (str), `genres` (str), `configs_count` (int), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `version` (object), `version.version_number` (str), `version.version_string` (str), `content_rating` (str), `content_rating_details` (str), `ranking` (int), `ranking_date` (datetime).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/rankings/{container}/{genre}/{country}). Verified 2026-09-09 signed-in documentation.

## get-appstore-sdks

`GET /appstore/sdks` — iOS SDKs » Directory

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `offset` | query | false | default=0; min=0 |
| `company.id` | query | false |  |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=rank |
| `sort_order` | query | false | default=desc |

Response fields: `rank` (int), `score` (float), `name` (str), `id` (int), `slug` (str), `total_installs` (int), `current_installs` (int), `uninstalls` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/sdks). Verified 2026-09-09 signed-in documentation.

## get-appstore-sdks-by-sdk-slug

`GET /appstore/sdks/{sdk_slug}` — iOS SDKs » Overview

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `sdk_slug` | path | true |  |

Response fields: `name` (str), `slug` (str), `description` (str), `tags` (list[str]), `url` (str), `logo_url` (str), `publisher_count` (int), `total_installs` (int), `current_installs` (int), `uninstalls` (int), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/sdks/{sdk_slug}). Verified 2026-09-09 signed-in documentation.

## get-appstore-sdks-by-sdk-slug-genres

`GET /appstore/sdks/{sdk_slug}/genres` — iOS SDKs » Genres

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `sdk_slug` | path | true |  |

Response fields: `genre` (str), `total_installs` (int), `current_installs` (int), `uninstalls` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/sdks/{sdk_slug}/genres). Verified 2026-09-09 signed-in documentation.

## get-appstore-sdks-by-sdk-slug-install-trend

`GET /appstore/sdks/{sdk_slug}/install-trend` — iOS SDKs » Install trend

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `sdk_slug` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |

Response fields: `total_installs` (int), `current_installs` (int), `uninstalls` (int), `date` (datetime).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/sdks/{sdk_slug}/install-trend). Verified 2026-09-09 signed-in documentation.

## get-appstore-sdks-by-sdk-slug-installs

`GET /appstore/sdks/{sdk_slug}/installs` — iOS SDKs » Installs

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `sdk_slug` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=rank |
| `sort_order` | query | false | default=desc |

Response fields: `id` (str), `bundle_id` (str), `title` (str), `rank` (int), `score` (float), `developer_id` (int), `developer` (str), `company_url` (str), `support_url` (str), `privacy_policy_url` (str), `seller` (str), `copyright` (str), `downloaded` (bool), `review_count` (int), `rating` (float), `rating_count` (int), `upload_date` (date), `release_date` (date), `small_artwork_url` (str), `large_artwork_url` (str), `xlarge_artwork_url` (str), `genre` (str), `genres` (str), `configs_count` (int), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `version` (object), `version.version_number` (str), `version.version_string` (str), `content_rating` (str), `content_rating_details` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/sdks/{sdk_slug}/installs). Verified 2026-09-09 signed-in documentation.

## get-appstore-sdks-by-sdk-slug-uninstalls

`GET /appstore/sdks/{sdk_slug}/uninstalls` — iOS SDKs » Uninstalls

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `sdk_slug` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=rank |
| `sort_order` | query | false | default=desc |

Response fields: `id` (str), `bundle_id` (str), `title` (str), `rank` (int), `score` (float), `developer_id` (int), `developer` (str), `company_url` (str), `support_url` (str), `privacy_policy_url` (str), `seller` (str), `copyright` (str), `downloaded` (bool), `review_count` (int), `rating` (float), `rating_count` (int), `upload_date` (date), `release_date` (date), `small_artwork_url` (str), `large_artwork_url` (str), `xlarge_artwork_url` (str), `genre` (str), `genres` (str), `configs_count` (int), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `version` (object), `version.version_number` (str), `version.version_string` (str), `content_rating` (str), `content_rating_details` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/sdks/{sdk_slug}/uninstalls). Verified 2026-09-09 signed-in documentation.

## get-appstore-sdks-by-sdk-slug-versions

`GET /appstore/sdks/{sdk_slug}/versions` — iOS SDKs » Versions

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `sdk_slug` | path | true |  |

Response fields: `id` (int), `version` (str), `name` (str), `summary` (str), `description` (str), `homepage` (str), `documentation_url` (str), `social_media_url` (str), `source` (object), `platform` (list[str]), `ios_deployment_target` (str), `authors` (object), `license` (object), `screenshots` (list[str]), `other_data` (object), `xcconfig` (object), `source_files` (list[str]), `public_header_files` (list[str]), `preserve_paths` (list[str]), `exclude_files` (list[str]), `vendored_libraries` (list[str]), `vendored_frameworks` (list[str]), `resources` (list[str]), `libraries` (list[str]), `frameworks` (list[str]), `weak_frameworks` (list[str]), `private_header_files` (list[str]), `requires_arc` (bool), `default_subspec` (str), `dependencies` (list[str]).

Sources: [provider contract](https://mixrank.com/api/documentation#/appstore/sdks/{sdk_slug}/versions). Verified 2026-09-09 signed-in documentation.

## get-companies

`GET /companies` — Companies » Directory

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=rank |
| `search` | query | false |  |

Response fields: `id` (int), `name` (str), `slug` (str), `rank` (int), `score` (float), `employees` (int), `rank_internetretailer` (int), `rank_incmagazine` (int), `rank_fortune` (int), `address` (object), `address.street_address` (str), `address.city` (str), `address.region` (str), `address.country` (str), `address.postal_code` (str), `linkedin` (object), `linkedin.id` (int), `linkedin.website` (str), `linkedin.url` (str), `linkedin.type` (str), `linkedin.size` (str), `linkedin.name` (str), `linkedin.founded` (int), `linkedin.description` (str), `linkedin.locations` (list[object]), `linkedin.locations[].address` (str), `linkedin.locations[].is_primary` (bool), `linkedin.country` (str), `linkedin.specialties` (list[str]), `linkedin.employee_count` (int), `linkedin.follower_count` (int), `linkedin.logo_url` (str), `linkedin.similar_pages` (list[int]), `web` (object), `web.site_count` (int), `web.tag_count` (int), `appstore` (object), `appstore.app_count` (int), `appstore.sdk_count` (int), `playstore` (object), `playstore.app_count` (int), `playstore.sdk_count` (int), `industries` (list[object]), `industries[].id` (int), `industries[].name` (str), `industries[].is_primary` (bool), `sic` (list[object]), `sic[].id` (str), `sic[].name` (str), `naics` (list[object]), `naics[].id` (str), `naics[].name` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/companies). Verified 2026-09-09 signed-in documentation.

## get-companies-by-company-id

`GET /companies/{company_id}` — Companies » Overview

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `company_id` | path | true |  |

Response fields: `id` (int), `name` (str), `slug` (str), `rank` (int), `score` (float), `employees` (int), `rank_internetretailer` (int), `rank_incmagazine` (int), `rank_fortune` (int), `address` (object), `address.street_address` (str), `address.city` (str), `address.region` (str), `address.country` (str), `address.postal_code` (str), `linkedin` (object), `linkedin.id` (int), `linkedin.website` (str), `linkedin.url` (str), `linkedin.type` (str), `linkedin.size` (str), `linkedin.name` (str), `linkedin.founded` (int), `linkedin.description` (str), `linkedin.locations` (list[object]), `linkedin.locations[].address` (str), `linkedin.locations[].is_primary` (bool), `linkedin.country` (str), `linkedin.specialties` (list[str]), `linkedin.employee_count` (int), `linkedin.follower_count` (int), `linkedin.logo_url` (str), `linkedin.similar_pages` (list[int]), `web` (object), `web.site_count` (int), `web.tag_count` (int), `appstore` (object), `appstore.app_count` (int), `appstore.sdk_count` (int), `playstore` (object), `playstore.app_count` (int), `playstore.sdk_count` (int), `industries` (list[object]), `industries[].id` (int), `industries[].name` (str), `industries[].is_primary` (bool), `sic` (list[object]), `sic[].id` (str), `sic[].name` (str), `naics` (list[object]), `naics[].id` (str), `naics[].name` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/companies/{company_id}). Verified 2026-09-09 signed-in documentation.

## get-companies-by-company-id-employee-metrics

`GET /companies/{company_id}/employee-metrics` — Companies » Employee Summary Metrics

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `company_id` | path | true |  |
| `job_tags` | query | false |  |

Response fields: `job_tag_label` (str), `job_tag_id` (int), `job_count_current` (int), `job_count_alltime` (int), `job_count_expired` (int), `job_first_seen` (date), `job_last_seen` (date).

Sources: [provider contract](https://mixrank.com/api/documentation#/companies/{company_id}/employee-metrics). Verified 2026-09-09 signed-in documentation.

## get-companies-by-company-id-employee-metrics-by-job-tag-id-timeseries

`GET /companies/{company_id}/employee-metrics/{job_tag_id}/timeseries` — Companies » Employee Trends

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `company_id` | path | true |  |
| `job_tag_id` | path | true |  |
| `since` | query | false |  |

Response fields: `job_tag_id` (int), `job_count_current` (int), `job_count_alltime` (int), `job_count_expired` (int), `month` (date).

Sources: [provider contract](https://mixrank.com/api/documentation#/companies/{company_id}/employee-metrics/{job_tag_id}/timeseries). Verified 2026-09-09 signed-in documentation.

## get-companies-by-company-id-timeseries

`GET /companies/{company_id}/timeseries` — Companies » Trends

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).

Sparse historical observations, with no uniform history guarantee. since is documented in the example rather than the parameter table.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `company_id` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=100; min=1; max=1000 |
| `since` | query | false |  |

Response fields: `date` (date), `employee_count` (int), `linkedin_follower_count` (int), `facebook_like_count` (int), `twitter_follower_count` (int), `seomoz_link_count` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/companies/{company_id}/timeseries). Verified 2026-09-09 signed-in documentation.

## get-companies-match

`GET /companies/match` — Companies » Match

Encoding: query. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).

Supports name, url, linkedin and linkedin_company_id. enable=company_financials is an optional billed add-on. An empty object is distinct from null.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `name` | query | false |  |
| `url` | query | false |  |
| `linkedin` | query | false |  |
| `linkedin_company_id` | query | false |  |
| `enable` | query | false |  |

Response fields: `id` (int), `name` (str), `slug` (str), `rank` (int), `score` (float), `employees` (int), `rank_internetretailer` (int), `rank_incmagazine` (int), `rank_fortune` (int), `address` (object), `address.street_address` (str), `address.city` (str), `address.region` (str), `address.country` (str), `address.postal_code` (str), `linkedin` (object), `linkedin.id` (int), `linkedin.website` (str), `linkedin.url` (str), `linkedin.type` (str), `linkedin.size` (str), `linkedin.name` (str), `linkedin.founded` (int), `linkedin.description` (str), `linkedin.locations` (list[object]), `linkedin.locations[].address` (str), `linkedin.locations[].is_primary` (bool), `linkedin.country` (str), `linkedin.specialties` (list[str]), `linkedin.employee_count` (int), `linkedin.follower_count` (int), `linkedin.logo_url` (str), `linkedin.similar_pages` (list[int]), `web` (object), `web.site_count` (int), `web.tag_count` (int), `appstore` (object), `appstore.app_count` (int), `appstore.sdk_count` (int), `playstore` (object), `playstore.app_count` (int), `playstore.sdk_count` (int), `industries` (list[object]), `industries[].id` (int), `industries[].name` (str), `industries[].is_primary` (bool), `sic` (list[object]), `sic[].id` (str), `sic[].name` (str), `naics` (list[object]), `naics[].id` (str), `naics[].name` (str), `company_financials` (object), `company_financials.funding` (object), `company_financials.funding.total_funding` (int), `company_financials.funding.total_funding_currency` (str), `company_financials.funding.round_count` (int), `company_financials.funding.rounds` (list[object]), `company_financials.funding.rounds[].round_type` (str), `company_financials.funding.rounds[].announced_date` (str), `company_financials.funding.rounds[].amount` (int), `company_financials.funding.rounds[].currency` (str), `company_financials.funding.rounds[].deal_status` (str), `company_financials.funding.rounds[].lead_investor` (str), `company_financials.funding.rounds[].investors` (list[str]).

Sources: [provider contract](https://mixrank.com/api/documentation#/companies/match). Verified 2026-09-09 signed-in documentation.

## get-companies-segment-preview

`GET /companies/segment-preview` — Companies » Company Segment Preview

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `query` | query | true |  |
| `limit` | query | false | default=100; min=1; max=100 |

Response fields: `total` (int), `is_truncated` (int), `results` (list[int]).

Sources: [provider contract](https://mixrank.com/api/documentation#/companies/segment-preview). Verified 2026-09-09 signed-in documentation.

## get-companies-segment-query

`GET /companies/segment-query` — Companies » Company Segment Query

Encoding: query. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `query` | query | true |  |

Response fields: `token` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/companies/segment-query). Verified 2026-09-09 signed-in documentation.

## get-echo

`GET /echo` — Account » Echo

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



Response fields: `status` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/echo). Verified 2026-09-09 signed-in documentation.

## get-email-prospect

`GET /email/prospect` — People » Get Prospect

Encoding: query. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).

May enrich contacts. webhook supports asynchronous 202 acceptance. Never automatically replay an uncertain enrichment request.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `full_name` | query | false |  |
| `first_name` | query | false |  |
| `last_name` | query | false |  |
| `person_id` | query | false |  |
| `domain` | query | false |  |
| `limit` | query | false | default=1; min=1; max=10 |
| `check_ownership` | query | false | default=t |
| `skip_sv` | query | false | default=f |
| `keep_ambiguous` | query | false | default=f |
| `downscore_edu` | query | false | default=t |
| `webhook` | query | false |  |

Response fields: `name` (str), `emails` (list[object]), `emails[].email` (str), `emails[].normalized` (str), `emails[].validity` (str), `emails[].validated_at` (datetime), `emails[].uses_greylist` (bool), `emails[].retry_after` (datetime), `emails[].gateway` (str), `emails[].is_disposable` (bool), `emails[].is_consumer` (bool), `emails[].is_domain_catchall` (bool), `emails[].mx_records` (list[str]).

Sources: [provider contract](https://mixrank.com/api/documentation#/email/prospect). Verified 2026-09-09 signed-in documentation.

## get-email-validate

`GET /email/validate` — Email » Validate

Encoding: query. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).

All toolkit calls are serialized per credential on one host. For lists use bulk jobs. A webhook returns 202 with an id; verification work continues after acceptance. Retry timing and validity remain evidence, not a guarantee.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `email` | query | true |  |
| `timeout` | query | false | default=120.0; min=0.0; max=120.0 |
| `maxage` | query | false |  |
| `strategy` | query | false | default=besteffort; enum=['cached', 'fetch', 'strict', 'besteffort'] |
| `private` | query | false | default=False |
| `webhook` | query | false |  |

Response fields: `email` (str), `normalized` (str), `validity` (str), `validated_at` (datetime), `uses_greylist` (bool), `retry_after` (datetime), `gateway` (str), `is_disposable` (bool), `is_consumer` (bool), `is_domain_catchall` (bool), `mx_records` (list[str]).

Sources: [provider contract](https://mixrank.com/api/documentation#/email/validate). Verified 2026-09-09 signed-in documentation.

## get-email-validate-bulk-job

`GET /email/validate/bulk-job` — Email » List Jobs

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).

POST uploads UTF-8 newline-separated emails as multipart file; limit 3000000 lines and 20 active jobs. private=true requires strategy=fetch and explicit deletion of downloaded output.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `page_size` | query | false | default=10; min=1; max=100 |
| `offset` | query | false | default=0; min=0 |

Response fields: `total` (int), `offset` (int), `page_size` (int), `results` (list[object]), `results[].id` (str), `results[].created_at` (datetime), `results[].lines_total` (int), `results[].lines_invalid` (int), `results[].lines_processed` (int), `results[].lines_refreshed` (int), `results[].strategy` (str), `results[].maxage` (datetime), `results[].completed_at` (datetime), `results[].download_url` (str), `results[].private` (bool).

Sources: [provider contract](https://mixrank.com/api/documentation#/email/validate/bulk-job). Verified 2026-09-09 signed-in documentation.

## get-email-validate-bulk-job-by-job-id

`GET /email/validate/bulk-job/{job_id}` — Email » Check Job

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `job_id` | path | true |  |

Response fields: `id` (str), `created_at` (datetime), `lines_total` (int), `lines_invalid` (int), `lines_processed` (int), `lines_refreshed` (int), `strategy` (str), `maxage` (datetime), `completed_at` (datetime), `download_url` (str), `private` (bool).

Sources: [provider contract](https://mixrank.com/api/documentation#/email/validate/bulk-job/{job_id}). Verified 2026-09-09 signed-in documentation.

## get-jobtags

`GET /jobtags` — Companies » Job Tags

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `tag_ids` | query | false |  |

Response fields: `id` (int), `label` (str), `parent_id` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/jobtags). Verified 2026-09-09 signed-in documentation.

## get-linkedin-company

`GET /linkedin/company` — LinkedIn » Company

Encoding: query. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).

Toolkit defaults to cached; fetch/strict/besteffort need explicit refresh. Synchronous LiveScan returns profile data; 512 is temporary fetch failure, 513 is permanent. No implicit strategy fallback.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `url` | query | false |  |
| `company.id` | query | false |  |
| `org_id` | query | false |  |
| `maxage` | query | false |  |
| `strategy` | query | false | default=cached; enum=['cached', 'fetch', 'strict', 'besteffort'] |

Response fields: `affiliated_pages` (list[int]), `company_headline` (str), `company_id` (int), `country_iso` (str), `country` (str), `country_name` (str), `created_at` (datetime), `crunchbase_funding` (list[object]), `crunchbase_funding[].round_name` (str), `crunchbase_funding[].round_date` (datetime), `crunchbase_funding[].round_amount` (str), `crunchbase_funding[].investor_names` (list[str]), `crunchbase_funding[].funding_round_count` (int), `crunchbase_funding[].investor_count` (int), `crunchbase_funding[].round_currency` (str), `crunchbase_funding[].crunchbase_company_name` (str), `crunchbase_funding[].funding_url` (str), `crunchbase_funding[].crunchbase_company_url` (str), `crunchbase_funding[].people_investors_urls` (list[str]), `crunchbase_funding[].organization_investors_urls` (list[str]), `description` (str), `employee_count` (int), `employees` (list[int]), `follower_count` (int), `founded_year` (int), `founded` (int), `hero_image` (str), `industries` (list[object]), `industries[].id` (int), `industries[].name` (int), `industries[].primary` (int), `inferred_location` (object), `linkedin_company_id` (int), `linkedin_url` (str), `url` (str), `locality` (str), `locations` (list[object]), `locations[].address` (str), `locations[].is_primary` (bool), `locations[].inferred_location` (object), `logo` (str), `logo_url` (str), `naics` (list[object]), `naics[].code` (int), `naics[].title` (str), `naics_codes` (list[object]), `naics_codes[].code` (int), `naics_codes[].title` (str), `name` (str), `org_id` (int), `postal_code` (str), `posts` (list[object]), `posts[].id` (int), `posts[].updated_at` (datetime), `posts[].posted_date_range` (object), `posts[].posted_date` (datetime), `posts[].content_html` (str), `posts[].likes_count` (int), `posts[].comments_count` (int), `posts[].image_url` (str), `posts[].video_url` (str), `posts[].external_url` (str), `posts[].embedded_post` (object), `region` (str), `similar_pages` (list[int]), `size` (str), `slug` (str), `slug_status` (str), `snap_id` (int), `specialties` (list[str]), `stock_exchange_code` (str), `street_address` (str), `subsidiary_linkedin_company_ids` (list[int]), `ticker` (str), `type` (str), `updated_at` (datetime), `last_refresh` (datetime), `website` (str), `phone_numbers` (str), `domain` (str), `industry` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/linkedin/company). Verified 2026-09-09 signed-in documentation.

## get-linkedin-company-posts

`GET /linkedin/company/posts` — LinkedIn » Company Posts

Encoding: query. Pagination: offset. Automatic replay eligible: false (refresh can further restrict this).

Toolkit defaults to cached; fetch/strict/besteffort need explicit refresh. Synchronous LiveScan returns profile data; 512 is temporary fetch failure, 513 is permanent. No implicit strategy fallback.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `sort_field` | query | false | default=id |
| `url` | query | false |  |
| `company.id` | query | false |  |
| `org_id` | query | false |  |
| `strategy` | query | false | default=cached; enum=['cached', 'fetch'] |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=10; min=1; max=100 |

Response fields: `id` (int), `linkedin_post_id` (int), `linkedin_org_id` (int), `url` (str), `post_date` (date), `content_html` (str), `likes_count` (int), `comments_count` (int), `image_urls` (str), `video_url` (str), `external_url` (str), `embedded_post` (str), `document_title` (str), `document_images` (list[object]), `is_repost` (bool), `reposts` (list[object]), `reposts[].reposter_linkedin_profile_id` (int), `reposts[].reposter_linkedin_company_id` (int), `reposts[].repost_data_activity_urn` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/linkedin/company/posts). Verified 2026-09-09 signed-in documentation.

## get-linkedin-job

`GET /linkedin/job` — LinkedIn » Job

Encoding: query. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).

Toolkit defaults to cached; fetch/strict/besteffort need explicit refresh. Synchronous LiveScan returns profile data; 512 is temporary fetch failure, 513 is permanent. No implicit strategy fallback.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `url` | query | false |  |
| `mixrank_id` | query | false |  |
| `linkedin_job_id` | query | false |  |
| `maxage` | query | false |  |
| `strategy` | query | false | default=cached; enum=['cached', 'fetch', 'strict', 'besteffort'] |

Response fields: `id` (int), `job_id` (int), `job_title` (str), `normalized_title` (object), `normalized_title.id` (int), `normalized_title.title` (str), `normalized_title.linkedin_title_id` (int), `job_posted_date` (datetime), `company` (object), `company.id` (int), `company.name` (str), `company.linkedin_org_id` (int), `company.universal_name` (str), `applicants` (str), `job_description` (str), `seniority` (object), `seniority.id` (int), `seniority.name` (str), `employment_type` (object), `employment_type.id` (int), `employment_type.name` (str), `job_functions` (list[object]), `job_functions[].id` (int), `job_functions[].name` (str), `industries` (list[object]), `industries[].id` (int), `industries[].name` (str), `salary` (object), `salary.currency` (str), `salary.unit` (str), `salary.min_salary` (int), `salary.max_salary` (int), `recruiter` (str), `benefits` (str), `address` (object), `address.street_address` (str), `address.locality` (str), `address.region` (str), `address.country` (str), `address.location` (str), `address.location_id` (str), `address.latitude` (str), `address.longitude` (str), `address.postal_code` (str), `applicant_count_history` (list[object]), `applicant_count_history[].date` (datetime), `applicant_count_history[].applicants` (int), `required_experience` (int), `required_qualification` (str), `job_application_url` (str), `job_application_id` (str), `job_close_date` (datetime), `url` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/linkedin/job). Verified 2026-09-09 signed-in documentation.

## get-linkedin-post-bulk-job

`GET /linkedin/post/bulk_job` — LinkedIn » Post Bulk Scan - List Jobs

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).

POST uploads UTF-8 post URLs as multipart file; limit 10000 lines and 20 active jobs. Submission cannot be cancelled. Use returned job id to inspect completion.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `page_size` | query | false | default=10; min=1; max=100 |
| `offset` | query | false | default=0; min=0 |

Response fields: `id` (str), `created_at` (datetime), `lines_total` (int), `lines_invalid` (int), `lines_processed` (int), `completed_at` (datetime).

Sources: [provider contract](https://mixrank.com/api/documentation#/linkedin/post/bulk_job). Verified 2026-09-09 signed-in documentation.

## get-linkedin-post-bulk-job-by-job-id

`GET /linkedin/post/bulk_job/{job_id}` — LinkedIn » Post Bulk Scan - Job Status

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `job_id` | path | true |  |

Response fields: `id` (str), `created_at` (datetime), `lines_total` (int), `lines_invalid` (int), `lines_processed` (int), `completed_at` (datetime).

Sources: [provider contract](https://mixrank.com/api/documentation#/linkedin/post/bulk_job/{job_id}). Verified 2026-09-09 signed-in documentation.

## get-linkedin-profile

`GET /linkedin/profile` — LinkedIn » Profile

Encoding: query. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).

Toolkit defaults to cached; fetch/strict/besteffort need explicit refresh. Synchronous LiveScan returns profile data; 512 is temporary fetch failure, 513 is permanent. No implicit strategy fallback.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `url` | query | false |  |
| `person.id` | query | false |  |
| `user_id` | query | false |  |
| `profile_id` | query | false |  |
| `maxage` | query | false |  |
| `strategy` | query | false | default=cached; enum=['cached', 'fetch', 'strict', 'besteffort'] |
| `experience_filter` | query | false | default=default; enum=['default', 'current', 'revisions', 'distincts'] |
| `education_filter` | query | false | default=default; enum=['default', 'current', 'revisions', 'distincts'] |
| `summary_filter` | query | false | default=default; enum=['default', 'last_complete'] |

Response fields: `person_id` (int), `profile_id` (int), `user_id` (int), `slug` (str), `url` (str), `name` (str), `first_name` (str), `last_name` (str), `company_name` (str), `org` (str), `title` (str), `country_name` (str), `country_iso` (str), `country` (str), `locality` (str), `location_name` (str), `headline` (str), `picture_url_orig` (str), `profile_pic` (str), `picture_url_copy` (str), `cover_image` (str), `connection_count` (int), `connections` (int), `recommender_count` (int), `num_recommenders` (int), `skills` (list[str]), `inferred_skills` (list[str]), `is_incomplete` (bool), `follower_count` (int), `num_followers` (int), `dob` (date), `jobs_count` (int), `updated_at` (datetime), `last_seen` (datetime), `last_refresh` (datetime), `summary` (str), `industry_id` (int), `industry_name` (str), `created_at` (datetime), `snap_id` (int), `privacy_redact` (bool), `slug_status` (str), `linkedin_company_id` (int), `public_experience_available` (bool), `public_education_available` (bool), `activity_at` (datetime), `is_memorial` (bool), `proserve_services` (object), `proserve_services.services` (object), `proserve_services.allow_free_message` (bool), `proserve_services.work_modality` (str), `proserve_services.work_location` (str), `slugs` (list[object]), `slugs[].slug` (object), `slugs[].status` (str), `specialties` (list[str]), `test_scores` (list[object]), `test_scores[].test_name` (str), `test_scores[].test_score` (str), `test_scores[].test_date` (datetime), `test_scores[].test_description` (str), `position` (object), `position.linkedin_company_id` (int), `position.company_name` (str), `position.title` (str), `position.summary` (str), `position.start_date_year` (int), `position.start_date_month` (int), `position.start_date` (date), `position.locality` (str), `interests` (list[str]), `inferred_location` (object), `inferred_location.latitude` (float), `inferred_location.longitude` (float), `inferred_location.locationstr` (str), `inferred_location.country_iso` (str), `inferred_location.address_line` (str), `inferred_location.admin_district` (str), `inferred_location.admin_district2` (str), `inferred_location.country_region` (str), `inferred_location.formatted_address` (str), `inferred_location.locality` (str), `inferred_location.postal_code` (str), `inferred_location.neighborhood` (str), `inferred_location.name` (str), `job_posts` (list[object]), `job_posts[].job_id` (int), `job_posts[].title` (str), `job_posts[].posted_timestamp` (datetime), `job_posts[].company_name` (str), `job_posts[].location` (str), `job_posts[].applicants` (int), `articles` (list[object]), `articles[].id` (int), `articles[].title` (str), `articles[].date` (date), `url_resources` (list[object]), `url_resources[].name` (str), `url_resources[].normalized_name` (str), `url_resources[].url` (str), `education` (list[object]), `education[].school_name` (str), `education[].degree` (str), `education[].grade` (str), `education[].start_date` (date), `education[].end_date` (date), `education[].activities` (str), `experience` (list[object]), `experience[].company` (str), `experience[].company_id` (int), `experience[].org_id` (int), `experience[].url` (int), `experience[].company_website` (str), `experience[].company_domain` (str), `experience[].title` (str), `experience[].summary` (str), `experience[].start_date` (date), `experience[].end_date` (date), `experience[].locality` (str), `experience[].is_current` (bool), `experience[].seniority` (list[object]), `experience[].job_function` (list[object]), `experience[].employment_type` (list[object]), `experience[].academic_qualification` (list[object]), `languages` (list[object]), `languages[].language` (str), `languages[].proficiency` (str), `recommendations` (list[object]), `recommendations[].recommender_linkedin_profile_id` (int), `recommendations[].recommender_linkedin_url` (str), `recommendations[].recommender_person_id` (int), `recommendations[].recommender_name` (str), `recommendations[].recommendation` (str), `patents` (list[object]), `patents[].title` (str), `patents[].country` (str), `patents[].number` (str), `patents[].description` (str), `patents[].url` (str), `publications` (list[object]), `publications[].title` (str), `publications[].publisher` (str), `publications[].date` (date), `publications[].summary` (str), `publications[].url` (str), `projects` (list[object]), `projects[].title` (str), `projects[].start_date` (str), `projects[].end_date` (str), `projects[].url` (str), `projects[].summary` (str), `volunteering` (list[object]), `volunteering[].company_name` (str), `volunteering[].company_id` (int), `volunteering[].company_website` (str), `volunteering[].company_domain` (str), `volunteering[].role` (str), `volunteering[].summary` (str), `volunteering[].cause` (str), `volunteering[].start_date` (str), `volunteering[].end_date` (str), `certifications` (list[object]), `certifications[].title` (str), `certifications[].summary` (str), `certifications[].verify_url` (str), `certifications[].credential_id` (str), `certifications[].company_name` (str), `certifications[].company_id` (int), `certifications[].date` (str), `certifications[].expire_date_month` (str), `certifications[].expire_date_year` (str), `awards` (list[object]), `awards[].title` (str), `awards[].summary` (str), `awards[].company_name` (str), `awards[].company_id` (int), `awards[].date` (str), `courses` (list[object]), `courses[].title` (str), `courses[].number` (str), `courses[].association` (str), `others_named` (list[object]), `others_named[].person_id` (int), `others_named[].url` (str), `people_also_viewed` (list[object]), `people_also_viewed[].person_id` (int), `people_also_viewed[].url` (str), `seniority` (list[object]), `seniority[].seniority` (str), `seniority[].id` (int), `job_function` (list[object]), `job_function[].job_function` (str), `job_function[].id` (int), `employment_type` (list[object]), `employment_type[].job_employment_type` (str), `employment_type[].id` (int), `academic_qualification` (list[object]), `academic_qualification[].academic_qualification` (str), `academic_qualification[].id` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/linkedin/profile). Verified 2026-09-09 signed-in documentation.

## get-linkedin-profile-bulk-job

`GET /linkedin/profile/bulk_job` — LinkedIn » List Jobs

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).

POST uploads UTF-8 profile URLs as multipart file; limit 250000 lines and 20 active jobs. Submitted jobs cannot be cancelled. completed_at and download_url indicate completion.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `page_size` | query | false | default=10; min=1; max=100 |
| `offset` | query | false | default=0; min=0 |

Response fields: `id` (str), `created_at` (datetime), `lines_total` (int), `lines_invalid` (int), `lines_processed` (int), `lines_refreshed` (int), `lines_cached` (int), `maxage` (int), `strategy` (str), `completed_at` (datetime), `download_url` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/linkedin/profile/bulk_job). Verified 2026-09-09 signed-in documentation.

## get-linkedin-profile-bulk-job-by-job-id

`GET /linkedin/profile/bulk_job/{job_id}` — LinkedIn » Check Bulk

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `job_id` | path | true |  |

Response fields: `id` (str), `created_at` (datetime), `lines_total` (int), `lines_invalid` (int), `lines_processed` (int), `lines_refreshed` (int), `lines_cached` (int), `maxage` (int), `strategy` (str), `completed_at` (datetime), `download_url` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/linkedin/profile/bulk_job/{job_id}). Verified 2026-09-09 signed-in documentation.

## get-linkedin-profile-posts

`GET /linkedin/profile/posts` — LinkedIn » Profile Posts

Encoding: query. Pagination: offset. Automatic replay eligible: false (refresh can further restrict this).

Toolkit defaults to cached; fetch/strict/besteffort need explicit refresh. Synchronous LiveScan returns profile data; 512 is temporary fetch failure, 513 is permanent. No implicit strategy fallback.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `sort_field` | query | false | default=id |
| `url` | query | false |  |
| `person.id` | query | false |  |
| `profile_id` | query | false |  |
| `strategy` | query | false | default=cached; enum=['cached', 'fetch'] |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=10; min=1; max=100 |

Response fields: `id` (int), `linkedin_post_id` (int), `url` (str), `post_date` (date), `content_html` (str), `likes_count` (int), `comments_count` (int), `image_urls` (str), `video_url` (str), `external_url` (str), `embedded_post` (str), `document_title` (str), `document_images` (list[object]), `is_repost` (bool), `reposts` (list[object]), `reposts[].reposter_linkedin_profile_id` (int), `reposts[].reposter_linkedin_company_id` (int), `reposts[].repost_data_activity_urn` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/linkedin/profile/posts). Verified 2026-09-09 signed-in documentation.

## get-person-by-id

`GET /person/{id}` — People » Overview

Encoding: query. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).

The documentation calls the path argument person_id; the URL template uses id. Contact add-ons use enable=b2b_emails,b2c_emails,edu_emails,directdials. Disabled data may be {} whereas absent data is null.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `id` | path | true |  |
| `enable` | query | false |  |
| `disable` | query | false |  |

Response fields: `id` (int), `name` (object), `name.first` (str), `name.middle` (str), `name.last` (str), `name.full` (str), `name.nickname` (str), `name.suffix` (str), `name.title` (str), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `company.rank` (int), `company.score` (float), `company.employees` (int), `company.rank_internetretailer` (int), `company.rank_incmagazine` (int), `company.rank_fortune` (int), `company.address` (object), `company.address.street_address` (str), `company.address.city` (str), `company.address.region` (str), `company.address.country` (str), `company.address.postal_code` (str), `company.linkedin` (object), `company.linkedin.id` (int), `company.linkedin.website` (str), `company.linkedin.url` (str), `company.linkedin.type` (str), `company.linkedin.size` (str), `company.linkedin.name` (str), `company.linkedin.founded` (int), `company.linkedin.description` (str), `company.linkedin.locations` (list[object]), `company.linkedin.locations[].address` (str), `company.linkedin.locations[].is_primary` (bool), `company.linkedin.country` (str), `company.linkedin.specialties` (list[str]), `company.linkedin.employee_count` (int), `company.linkedin.follower_count` (int), `company.linkedin.logo_url` (str), `company.linkedin.similar_pages` (list[int]), `company.web` (object), `company.web.site_count` (int), `company.web.tag_count` (int), `company.appstore` (object), `company.appstore.app_count` (int), `company.appstore.sdk_count` (int), `company.playstore` (object), `company.playstore.app_count` (int), `company.playstore.sdk_count` (int), `company.industries` (list[object]), `company.industries[].id` (int), `company.industries[].name` (str), `company.industries[].is_primary` (bool), `company.sic` (list[object]), `company.sic[].id` (str), `company.sic[].name` (str), `company.naics` (list[object]), `company.naics[].id` (str), `company.naics[].name` (str), `linkedin` (object), `linkedin.id` (int), `linkedin.profile_pic` (str), `linkedin.connections` (int), `linkedin.org` (str), `linkedin.name` (object), `linkedin.name.first` (str), `linkedin.name.middle` (str), `linkedin.name.last` (str), `linkedin.name.full` (str), `linkedin.name.nickname` (str), `linkedin.name.suffix` (str), `linkedin.name.title` (str), `linkedin.title` (str), `linkedin.headline` (str), `linkedin.country` (str), `linkedin.industry` (object), `linkedin.industry.id` (int), `linkedin.industry.name` (str), `linkedin.num_recommenders` (int), `linkedin.summary` (str), `linkedin.url` (str), `linkedin.dob` (date), `linkedin.slug_status` (str), `linkedin.skills` (list[str]), `linkedin.location` (object), `linkedin.location.text` (str), `linkedin.location.region` (str), `linkedin.location.city` (str), `linkedin.location.district` (str), `linkedin.location.country` (object), `linkedin.location.country.id` (int), `linkedin.location.country.name` (str), `linkedin.positions` (list[object]), `linkedin.positions[].id` (int), `linkedin.positions[].company_id` (int), `linkedin.positions[].linkedin_company_id` (int), `linkedin.positions[].company_url` (str), `linkedin.positions[].company_name` (str), `linkedin.positions[].start_date` (date), `linkedin.positions[].start_date_year` (int), `linkedin.positions[].start_date_month` (int), `linkedin.positions[].end_date` (date), `linkedin.positions[].end_date_year` (int), `linkedin.positions[].end_date_month` (int), `linkedin.positions[].title` (str), `linkedin.positions[].locality` (str), `linkedin.positions[].is_current` (bool), `linkedin.education` (list[object]), `linkedin.education[].school_name` (str), `linkedin.education[].field_of_study` (int), `linkedin.education[].degree` (str), `linkedin.education[].grade` (str), `linkedin.education[].activities` (str), `linkedin.education[].notes` (str), `linkedin.education[].start_date` (date), `linkedin.education[].start_date_year` (int), `linkedin.education[].start_date_month` (int), `linkedin.education[].end_date` (date), `linkedin.education[].end_date_year` (int), `linkedin.education[].end_date_month` (int), `linkedin.volunteering` (list[object]), `linkedin.volunteering[].company_name` (str), `linkedin.volunteering[].company_id` (int), `linkedin.volunteering[].role` (str), `linkedin.volunteering[].summary` (str), `linkedin.volunteering[].cause` (str), `linkedin.volunteering[].start_date` (date), `linkedin.volunteering[].start_date_year` (int), `linkedin.volunteering[].start_date_month` (int), `linkedin.volunteering[].end_date` (date), `linkedin.volunteering[].end_date_year` (int), `linkedin.volunteering[].end_date_month` (int), `linkedin.languages` (list[object]), `linkedin.languages[].language` (str), `linkedin.languages[].proficiency` (str), `linkedin.certifications` (list[object]), `linkedin.certifications[].title` (str), `linkedin.certifications[].summary` (str), `linkedin.certifications[].verify_url` (str), `linkedin.certifications[].credential_id` (str), `linkedin.certifications[].company_name` (str), `linkedin.certifications[].company_id` (int), `linkedin.certifications[].date_year` (str), `linkedin.certifications[].date_month` (str), `linkedin.projects` (list[object]), `linkedin.projects[].title` (str), `linkedin.projects[].start_date` (str), `linkedin.projects[].end_date` (str), `linkedin.projects[].url` (str), `linkedin.recommendations` (list[object]), `linkedin.recommendations[].recommender_linkedin_profile_id` (int), `linkedin.recommendations[].recommender_profile_url` (int), `linkedin.recommendations[].recommender_person_id` (int), `linkedin.recommendations[].recommender_name` (str), `linkedin.recommendations[].recommendation` (str), `linkedin.updated_date` (date), `linkedin.seniority` (object), `linkedin.seniority.id` (int), `linkedin.seniority.seniority` (str), `linkedin.job_function` (object), `linkedin.job_function.id` (int), `linkedin.job_function.job_function` (str), `linkedin.employment_type` (object), `linkedin.employment_type.id` (int), `linkedin.employment_type.job_employment_type` (str), `linkedin.academic_qualification` (object), `linkedin.academic_qualification.id` (int), `linkedin.academic_qualification.academic_qualification` (str), `indeed` (object), `indeed.id` (str), `indeed.summary` (str), `indeed.vanity_name` (str), `indeed.name` (object), `indeed.name.first` (str), `indeed.name.middle` (str), `indeed.name.last` (str), `indeed.name.full` (str), `indeed.name.nickname` (str), `indeed.name.suffix` (str), `indeed.name.title` (str), `indeed.title` (str), `indeed.headline` (str), `indeed.url` (str), `indeed.edited_date` (date), `indeed.employment_eligibility` (str), `indeed.relocation_status` (str), `indeed.additional_info` (str), `indeed.locality` (str), `indeed.city` (str), `indeed.state` (str), `indeed.positions` (list[object]), `indeed.positions[].id` (int), `indeed.positions[].company_id` (int), `indeed.positions[].linkedin_company_id` (int), `indeed.positions[].company_url` (str), `indeed.positions[].company_name` (str), `indeed.positions[].start_date` (date), `indeed.positions[].start_date_year` (int), `indeed.positions[].start_date_month` (int), `indeed.positions[].end_date` (date), `indeed.positions[].end_date_year` (int), `indeed.positions[].end_date_month` (int), `indeed.positions[].title` (str), `indeed.positions[].locality` (str), `indeed.positions[].is_current` (bool), `directdials` (list[str]), `b2b_emails` (list[object]), `b2b_emails[].email` (str), `b2b_emails[].domain` (str), `b2c_emails` (list[object]), `b2c_emails[].email` (str), `b2c_emails[].domain` (str), `edu_emails` (list[object]), `edu_emails[].email` (str), `edu_emails[].domain` (str), `twitter` (object), `twitter.id` (int), `twitter.name` (object), `twitter.name.first` (str), `twitter.name.middle` (str), `twitter.name.last` (str), `twitter.name.full` (str), `twitter.name.nickname` (str), `twitter.name.suffix` (str), `twitter.name.title` (str), `twitter.url` (str), `twitter.signup_at` (date), `twitter.slug` (str), `twitter.follower_count` (int), `github` (object), `github.id` (int), `github.user_id` (int), `github.name` (object), `github.name.first` (str), `github.name.middle` (str), `github.name.last` (str), `github.name.full` (str), `github.name.nickname` (str), `github.name.suffix` (str), `github.name.title` (str), `github.login` (str), `github.url` (str), `github.updated_at` (date), `github.added_at` (date), `github.modified_at` (date), `github.site_admin` (bool), `github.avatar_url` (bool), `github.public_repos` (int), `github.public_gists` (int), `github.followers` (int), `github.followings` (int), `gplus` (object), `gplus.id` (str), `gplus.name` (object), `gplus.name.first` (str), `gplus.name.middle` (str), `gplus.name.last` (str), `gplus.name.full` (str), `gplus.name.nickname` (str), `gplus.name.suffix` (str), `gplus.name.title` (str), `gplus.headline` (str), `gplus.headline2` (str), `gplus.url` (str), `gplus.picture_url` (str), `gplus.company` (str), `gplus.school` (str), `gplus.location` (str), `facebook` (object), `facebook.id` (int), `facebook.username` (str), `facebook.user_id` (int), `facebook.url` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/person/{id}). Verified 2026-09-09 signed-in documentation.

## get-person-match

`GET /person/match` — People » Match

Encoding: query. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).

Supply email, phone, social_url, or name plus company/domain evidence. All add-ons default off. page_size defaults to 1 at the provider; request more candidates to assess ambiguity. Preserve confidence and conflicting evidence.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `email` | query | false |  |
| `phone` | query | false |  |
| `domain` | query | false |  |
| `company_name` | query | false |  |
| `name` | query | false |  |
| `first_name` | query | false |  |
| `last_name` | query | false |  |
| `page_size` | query | false | default=1; min=1; max=100 |
| `fuzzy` | query | false | default=t |
| `enable` | query | false |  |
| `disable` | query | false |  |
| `social_url` | query | false |  |
| `directdials_limit` | query | false |  |
| `b2b_emails_limit` | query | false |  |
| `b2c_emails_limit` | query | false |  |
| `edu_emails_limit` | query | false |  |

Response fields: `id` (int), `name` (object), `name.first` (str), `name.middle` (str), `name.last` (str), `name.full` (str), `name.nickname` (str), `name.suffix` (str), `name.title` (str), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `company.rank` (int), `company.score` (float), `company.employees` (int), `company.rank_internetretailer` (int), `company.rank_incmagazine` (int), `company.rank_fortune` (int), `company.address` (object), `company.address.street_address` (str), `company.address.city` (str), `company.address.region` (str), `company.address.country` (str), `company.address.postal_code` (str), `company.linkedin` (object), `company.linkedin.id` (int), `company.linkedin.website` (str), `company.linkedin.url` (str), `company.linkedin.type` (str), `company.linkedin.size` (str), `company.linkedin.name` (str), `company.linkedin.founded` (int), `company.linkedin.description` (str), `company.linkedin.locations` (list[object]), `company.linkedin.locations[].address` (str), `company.linkedin.locations[].is_primary` (bool), `company.linkedin.country` (str), `company.linkedin.specialties` (list[str]), `company.linkedin.employee_count` (int), `company.linkedin.follower_count` (int), `company.linkedin.logo_url` (str), `company.linkedin.similar_pages` (list[int]), `company.web` (object), `company.web.site_count` (int), `company.web.tag_count` (int), `company.appstore` (object), `company.appstore.app_count` (int), `company.appstore.sdk_count` (int), `company.playstore` (object), `company.playstore.app_count` (int), `company.playstore.sdk_count` (int), `company.industries` (list[object]), `company.industries[].id` (int), `company.industries[].name` (str), `company.industries[].is_primary` (bool), `company.sic` (list[object]), `company.sic[].id` (str), `company.sic[].name` (str), `company.naics` (list[object]), `company.naics[].id` (str), `company.naics[].name` (str), `linkedin` (object), `linkedin.id` (int), `linkedin.profile_pic` (str), `linkedin.connections` (int), `linkedin.org` (str), `linkedin.name` (object), `linkedin.name.first` (str), `linkedin.name.middle` (str), `linkedin.name.last` (str), `linkedin.name.full` (str), `linkedin.name.nickname` (str), `linkedin.name.suffix` (str), `linkedin.name.title` (str), `linkedin.title` (str), `linkedin.headline` (str), `linkedin.country` (str), `linkedin.industry` (object), `linkedin.industry.id` (int), `linkedin.industry.name` (str), `linkedin.num_recommenders` (int), `linkedin.summary` (str), `linkedin.url` (str), `linkedin.dob` (date), `linkedin.slug_status` (str), `linkedin.skills` (list[str]), `linkedin.location` (object), `linkedin.location.text` (str), `linkedin.location.region` (str), `linkedin.location.city` (str), `linkedin.location.district` (str), `linkedin.location.country` (object), `linkedin.location.country.id` (int), `linkedin.location.country.name` (str), `linkedin.positions` (list[object]), `linkedin.positions[].id` (int), `linkedin.positions[].company_id` (int), `linkedin.positions[].linkedin_company_id` (int), `linkedin.positions[].company_url` (str), `linkedin.positions[].company_name` (str), `linkedin.positions[].start_date` (date), `linkedin.positions[].start_date_year` (int), `linkedin.positions[].start_date_month` (int), `linkedin.positions[].end_date` (date), `linkedin.positions[].end_date_year` (int), `linkedin.positions[].end_date_month` (int), `linkedin.positions[].title` (str), `linkedin.positions[].locality` (str), `linkedin.positions[].is_current` (bool), `linkedin.education` (list[object]), `linkedin.education[].school_name` (str), `linkedin.education[].field_of_study` (int), `linkedin.education[].degree` (str), `linkedin.education[].grade` (str), `linkedin.education[].activities` (str), `linkedin.education[].notes` (str), `linkedin.education[].start_date` (date), `linkedin.education[].start_date_year` (int), `linkedin.education[].start_date_month` (int), `linkedin.education[].end_date` (date), `linkedin.education[].end_date_year` (int), `linkedin.education[].end_date_month` (int), `linkedin.volunteering` (list[object]), `linkedin.volunteering[].company_name` (str), `linkedin.volunteering[].company_id` (int), `linkedin.volunteering[].role` (str), `linkedin.volunteering[].summary` (str), `linkedin.volunteering[].cause` (str), `linkedin.volunteering[].start_date` (date), `linkedin.volunteering[].start_date_year` (int), `linkedin.volunteering[].start_date_month` (int), `linkedin.volunteering[].end_date` (date), `linkedin.volunteering[].end_date_year` (int), `linkedin.volunteering[].end_date_month` (int), `linkedin.languages` (list[object]), `linkedin.languages[].language` (str), `linkedin.languages[].proficiency` (str), `linkedin.certifications` (list[object]), `linkedin.certifications[].title` (str), `linkedin.certifications[].summary` (str), `linkedin.certifications[].verify_url` (str), `linkedin.certifications[].credential_id` (str), `linkedin.certifications[].company_name` (str), `linkedin.certifications[].company_id` (int), `linkedin.certifications[].date_year` (str), `linkedin.certifications[].date_month` (str), `linkedin.projects` (list[object]), `linkedin.projects[].title` (str), `linkedin.projects[].start_date` (str), `linkedin.projects[].end_date` (str), `linkedin.projects[].url` (str), `linkedin.recommendations` (list[object]), `linkedin.recommendations[].recommender_linkedin_profile_id` (int), `linkedin.recommendations[].recommender_profile_url` (int), `linkedin.recommendations[].recommender_person_id` (int), `linkedin.recommendations[].recommender_name` (str), `linkedin.recommendations[].recommendation` (str), `linkedin.updated_date` (date), `linkedin.seniority` (object), `linkedin.seniority.id` (int), `linkedin.seniority.seniority` (str), `linkedin.job_function` (object), `linkedin.job_function.id` (int), `linkedin.job_function.job_function` (str), `linkedin.employment_type` (object), `linkedin.employment_type.id` (int), `linkedin.employment_type.job_employment_type` (str), `linkedin.academic_qualification` (object), `linkedin.academic_qualification.id` (int), `linkedin.academic_qualification.academic_qualification` (str), `indeed` (object), `indeed.id` (str), `indeed.summary` (str), `indeed.vanity_name` (str), `indeed.name` (object), `indeed.name.first` (str), `indeed.name.middle` (str), `indeed.name.last` (str), `indeed.name.full` (str), `indeed.name.nickname` (str), `indeed.name.suffix` (str), `indeed.name.title` (str), `indeed.title` (str), `indeed.headline` (str), `indeed.url` (str), `indeed.edited_date` (date), `indeed.employment_eligibility` (str), `indeed.relocation_status` (str), `indeed.additional_info` (str), `indeed.locality` (str), `indeed.city` (str), `indeed.state` (str), `indeed.positions` (list[object]), `indeed.positions[].id` (int), `indeed.positions[].company_id` (int), `indeed.positions[].linkedin_company_id` (int), `indeed.positions[].company_url` (str), `indeed.positions[].company_name` (str), `indeed.positions[].start_date` (date), `indeed.positions[].start_date_year` (int), `indeed.positions[].start_date_month` (int), `indeed.positions[].end_date` (date), `indeed.positions[].end_date_year` (int), `indeed.positions[].end_date_month` (int), `indeed.positions[].title` (str), `indeed.positions[].locality` (str), `indeed.positions[].is_current` (bool), `directdials` (list[str]), `b2b_emails` (list[object]), `b2b_emails[].email` (str), `b2b_emails[].domain` (str), `b2c_emails` (list[object]), `b2c_emails[].email` (str), `b2c_emails[].domain` (str), `edu_emails` (list[object]), `edu_emails[].email` (str), `edu_emails[].domain` (str), `twitter` (object), `twitter.id` (int), `twitter.name` (object), `twitter.name.first` (str), `twitter.name.middle` (str), `twitter.name.last` (str), `twitter.name.full` (str), `twitter.name.nickname` (str), `twitter.name.suffix` (str), `twitter.name.title` (str), `twitter.url` (str), `twitter.signup_at` (date), `twitter.slug` (str), `twitter.follower_count` (int), `github` (object), `github.id` (int), `github.user_id` (int), `github.name` (object), `github.name.first` (str), `github.name.middle` (str), `github.name.last` (str), `github.name.full` (str), `github.name.nickname` (str), `github.name.suffix` (str), `github.name.title` (str), `github.login` (str), `github.url` (str), `github.updated_at` (date), `github.added_at` (date), `github.modified_at` (date), `github.site_admin` (bool), `github.avatar_url` (bool), `github.public_repos` (int), `github.public_gists` (int), `github.followers` (int), `github.followings` (int), `gplus` (object), `gplus.id` (str), `gplus.name` (object), `gplus.name.first` (str), `gplus.name.middle` (str), `gplus.name.last` (str), `gplus.name.full` (str), `gplus.name.nickname` (str), `gplus.name.suffix` (str), `gplus.name.title` (str), `gplus.headline` (str), `gplus.headline2` (str), `gplus.url` (str), `gplus.picture_url` (str), `gplus.company` (str), `gplus.school` (str), `gplus.location` (str), `facebook` (object), `facebook.id` (int), `facebook.username` (str), `facebook.user_id` (int), `facebook.url` (str), `match` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/person/match). Verified 2026-09-09 signed-in documentation.

## get-person-search

`GET /person/search` — People » Search

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `search` | query | true |  |
| `page_size` | query | false | default=4; min=1; max=100 |

Response fields: `id` (int), `name` (object), `name.first` (str), `name.middle` (str), `name.last` (str), `name.full` (str), `name.nickname` (str), `name.suffix` (str), `name.title` (str), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `company.rank` (int), `company.score` (float), `company.employees` (int), `company.rank_internetretailer` (int), `company.rank_incmagazine` (int), `company.rank_fortune` (int), `company.address` (object), `company.address.street_address` (str), `company.address.city` (str), `company.address.region` (str), `company.address.country` (str), `company.address.postal_code` (str), `company.linkedin` (object), `company.linkedin.id` (int), `company.linkedin.website` (str), `company.linkedin.url` (str), `company.linkedin.type` (str), `company.linkedin.size` (str), `company.linkedin.name` (str), `company.linkedin.founded` (int), `company.linkedin.description` (str), `company.linkedin.locations` (list[object]), `company.linkedin.locations[].address` (str), `company.linkedin.locations[].is_primary` (bool), `company.linkedin.country` (str), `company.linkedin.specialties` (list[str]), `company.linkedin.employee_count` (int), `company.linkedin.follower_count` (int), `company.linkedin.logo_url` (str), `company.linkedin.similar_pages` (list[int]), `company.web` (object), `company.web.site_count` (int), `company.web.tag_count` (int), `company.appstore` (object), `company.appstore.app_count` (int), `company.appstore.sdk_count` (int), `company.playstore` (object), `company.playstore.app_count` (int), `company.playstore.sdk_count` (int), `company.industries` (list[object]), `company.industries[].id` (int), `company.industries[].name` (str), `company.industries[].is_primary` (bool), `company.sic` (list[object]), `company.sic[].id` (str), `company.sic[].name` (str), `company.naics` (list[object]), `company.naics[].id` (str), `company.naics[].name` (str), `linkedin` (object), `linkedin.id` (int), `linkedin.profile_pic` (str), `linkedin.connections` (int), `linkedin.org` (str), `linkedin.name` (object), `linkedin.name.first` (str), `linkedin.name.middle` (str), `linkedin.name.last` (str), `linkedin.name.full` (str), `linkedin.name.nickname` (str), `linkedin.name.suffix` (str), `linkedin.name.title` (str), `linkedin.title` (str), `linkedin.headline` (str), `linkedin.country` (str), `linkedin.industry` (object), `linkedin.industry.id` (int), `linkedin.industry.name` (str), `linkedin.num_recommenders` (int), `linkedin.summary` (str), `linkedin.url` (str), `linkedin.dob` (date), `linkedin.slug_status` (str), `linkedin.skills` (list[str]), `linkedin.location` (object), `linkedin.location.text` (str), `linkedin.location.region` (str), `linkedin.location.city` (str), `linkedin.location.district` (str), `linkedin.location.country` (object), `linkedin.location.country.id` (int), `linkedin.location.country.name` (str), `linkedin.positions` (list[object]), `linkedin.positions[].id` (int), `linkedin.positions[].company_id` (int), `linkedin.positions[].linkedin_company_id` (int), `linkedin.positions[].company_url` (str), `linkedin.positions[].company_name` (str), `linkedin.positions[].start_date` (date), `linkedin.positions[].start_date_year` (int), `linkedin.positions[].start_date_month` (int), `linkedin.positions[].end_date` (date), `linkedin.positions[].end_date_year` (int), `linkedin.positions[].end_date_month` (int), `linkedin.positions[].title` (str), `linkedin.positions[].locality` (str), `linkedin.positions[].is_current` (bool), `linkedin.education` (list[object]), `linkedin.education[].school_name` (str), `linkedin.education[].field_of_study` (int), `linkedin.education[].degree` (str), `linkedin.education[].grade` (str), `linkedin.education[].activities` (str), `linkedin.education[].notes` (str), `linkedin.education[].start_date` (date), `linkedin.education[].start_date_year` (int), `linkedin.education[].start_date_month` (int), `linkedin.education[].end_date` (date), `linkedin.education[].end_date_year` (int), `linkedin.education[].end_date_month` (int), `linkedin.volunteering` (list[object]), `linkedin.volunteering[].company_name` (str), `linkedin.volunteering[].company_id` (int), `linkedin.volunteering[].role` (str), `linkedin.volunteering[].summary` (str), `linkedin.volunteering[].cause` (str), `linkedin.volunteering[].start_date` (date), `linkedin.volunteering[].start_date_year` (int), `linkedin.volunteering[].start_date_month` (int), `linkedin.volunteering[].end_date` (date), `linkedin.volunteering[].end_date_year` (int), `linkedin.volunteering[].end_date_month` (int), `linkedin.languages` (list[object]), `linkedin.languages[].language` (str), `linkedin.languages[].proficiency` (str), `linkedin.certifications` (list[object]), `linkedin.certifications[].title` (str), `linkedin.certifications[].summary` (str), `linkedin.certifications[].verify_url` (str), `linkedin.certifications[].credential_id` (str), `linkedin.certifications[].company_name` (str), `linkedin.certifications[].company_id` (int), `linkedin.certifications[].date_year` (str), `linkedin.certifications[].date_month` (str), `linkedin.projects` (list[object]), `linkedin.projects[].title` (str), `linkedin.projects[].start_date` (str), `linkedin.projects[].end_date` (str), `linkedin.projects[].url` (str), `linkedin.recommendations` (list[object]), `linkedin.recommendations[].recommender_linkedin_profile_id` (int), `linkedin.recommendations[].recommender_profile_url` (int), `linkedin.recommendations[].recommender_person_id` (int), `linkedin.recommendations[].recommender_name` (str), `linkedin.recommendations[].recommendation` (str), `linkedin.updated_date` (date), `linkedin.seniority` (object), `linkedin.seniority.id` (int), `linkedin.seniority.seniority` (str), `linkedin.job_function` (object), `linkedin.job_function.id` (int), `linkedin.job_function.job_function` (str), `linkedin.employment_type` (object), `linkedin.employment_type.id` (int), `linkedin.employment_type.job_employment_type` (str), `linkedin.academic_qualification` (object), `linkedin.academic_qualification.id` (int), `linkedin.academic_qualification.academic_qualification` (str), `indeed` (object), `indeed.id` (str), `indeed.summary` (str), `indeed.vanity_name` (str), `indeed.name` (object), `indeed.name.first` (str), `indeed.name.middle` (str), `indeed.name.last` (str), `indeed.name.full` (str), `indeed.name.nickname` (str), `indeed.name.suffix` (str), `indeed.name.title` (str), `indeed.title` (str), `indeed.headline` (str), `indeed.url` (str), `indeed.edited_date` (date), `indeed.employment_eligibility` (str), `indeed.relocation_status` (str), `indeed.additional_info` (str), `indeed.locality` (str), `indeed.city` (str), `indeed.state` (str), `indeed.positions` (list[object]), `indeed.positions[].id` (int), `indeed.positions[].company_id` (int), `indeed.positions[].linkedin_company_id` (int), `indeed.positions[].company_url` (str), `indeed.positions[].company_name` (str), `indeed.positions[].start_date` (date), `indeed.positions[].start_date_year` (int), `indeed.positions[].start_date_month` (int), `indeed.positions[].end_date` (date), `indeed.positions[].end_date_year` (int), `indeed.positions[].end_date_month` (int), `indeed.positions[].title` (str), `indeed.positions[].locality` (str), `indeed.positions[].is_current` (bool).

Sources: [provider contract](https://mixrank.com/api/documentation#/person/search). Verified 2026-09-09 signed-in documentation.

## get-person-segment-preview

`GET /person/segment-preview` — People » Person Segment Preview

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `query` | query | true |  |
| `limit` | query | false | default=100; min=1; max=100 |

Response fields: `total` (int), `is_truncated` (int), `results` (list[int]).

Sources: [provider contract](https://mixrank.com/api/documentation#/person/segment-preview). Verified 2026-09-09 signed-in documentation.

## get-person-segment-query

`GET /person/segment-query` — People » Person Segment Query

Encoding: query. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `query` | query | true |  |

Response fields: `token` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/person/segment-query). Verified 2026-09-09 signed-in documentation.

## get-person-social-url-lookup

`GET /person/social-url-lookup` — People » Social URL Lookup

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `url` | query | false |  |

Response fields: `person` (object), `person.id` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/person/social-url-lookup). Verified 2026-09-09 signed-in documentation.

## get-playstore-apks-by-package-name-by-version-id

`GET /playstore/apks/{package_name}/{version_id}` — Play Store Apps » APK

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `package_name` | path | true |  |
| `version_id` | path | true |  |

Response fields: `package_name` (str), `version_id` (int), `error` (str), `version_string` (str), `major_version_number` (int), `release_date` (date), `changes_html` (str), `namespace_count` (int), `permission_count` (int), `splits` (list[object]), `splits[].name` (str), `splits[].file_type` (int), `splits[].file_extension` (str), `splits[].byte_size` (int), `splits[].downloaded_at` (date), `splits[].last_modified` (date), `splits[].md5` (str), `splits[].sha256` (str), `splits[].sha1` (str), `splits[].namespace_count` (int), `splits[].permission_count` (int), `splits[].download_link` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/apks/{package_name}/{version_id}). Verified 2026-09-09 signed-in documentation.

## get-playstore-apps

`GET /playstore/apps` — Play Store Apps » Directory

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=downloads |
| `company.id` | query | false |  |
| `search` | query | false |  |

Response fields: `pname` (str), `title` (str), `rank` (int), `score` (float), `developer_id` (int), `developer_name` (str), `developer_email` (str), `developer_website` (str), `downloads` (int), `removed` (bool), `review_count` (int), `rating` (float), `rating_count` (int), `upload_date` (date), `release_date` (date), `small_icon_url` (str), `large_icon_url` (str), `category` (str), `categories` (str), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `version` (object), `version.major_version_number` (int), `version.version_code` (int), `version.version_string` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/apps). Verified 2026-09-09 signed-in documentation.

## get-playstore-apps-by-pname

`GET /playstore/apps/{pname}` — Play Store Apps » Overview

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `pname` | path | true |  |

Response fields: `pname` (str), `title` (str), `rank` (int), `score` (float), `developer_id` (int), `developer_name` (str), `developer_email` (str), `developer_website` (str), `downloads` (int), `removed` (bool), `review_count` (int), `rating` (float), `rating_count` (int), `upload_date` (date), `release_date` (date), `small_icon_url` (str), `large_icon_url` (str), `category` (str), `categories` (str), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `version` (object), `version.major_version_number` (int), `version.version_code` (int), `version.version_string` (str), `permissions` (list[str]), `ratings` (object), `ratings.1` (int), `ratings.2` (int), `ratings.3` (int), `ratings.4` (int), `ratings.5` (int), `recent_changes_html` (str), `description_html` (str), `privacy_policy` (str), `has_iap` (bool), `iap_description` (str), `contains_ads` (bool), `editors_choice` (bool), `installation_size` (int), `content_rating` (int), `content_rating_details` (list[str]), `version` (object), `version.major_version_number` (int), `version.version_code` (int), `version.version_string` (str), `prices` (list[object]), `prices[].native_currency` (str), `prices[].native_price` (int), `prices[].converted_currency` (str), `prices[].converted_price` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/apps/{pname}). Verified 2026-09-09 signed-in documentation.

## get-playstore-apps-by-pname-images

`GET /playstore/apps/{pname}/images` — Play Store Apps » Images

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `pname` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=100; min=1; max=100 |

Response fields: `url` (str), `type` (object).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/apps/{pname}/images). Verified 2026-09-09 signed-in documentation.

## get-playstore-apps-by-pname-namespaces

`GET /playstore/apps/{pname}/namespaces` — Play Store Apps » Namespaces

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `pname` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=version_code_min |
| `sort_order` | query | false | default=desc |

Response fields: `namespace` (str), `upload_date_min` (date), `upload_date_max` (date), `version_code_min` (int), `version_code_max` (int), `currently_installed` (bool).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/apps/{pname}/namespaces). Verified 2026-09-09 signed-in documentation.

## get-playstore-apps-by-pname-rankings

`GET /playstore/apps/{pname}/rankings` — Play Store Apps » Rankings Summary

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `pname` | path | true |  |
| `country` | query | false | default=US |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=current_rank |
| `sort_order` | query | false | default=desc |

Response fields: `container` (str), `country` (str), `category` (str), `current_rank` (int), `highest_rank` (int), `first_ranked` (date), `last_ranked` (date).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/apps/{pname}/rankings). Verified 2026-09-09 signed-in documentation.

## get-playstore-apps-by-pname-rankings-history

`GET /playstore/apps/{pname}/rankings-history` — Play Store Apps » Rankings History

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `pname` | path | true |  |
| `country` | query | false | default=US |
| `container` | query | true |  |
| `category` | query | false |  |
| `before` | query | false |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |

Response fields: `country` (str), `ranking` (int), `date` (datetime).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/apps/{pname}/rankings-history). Verified 2026-09-09 signed-in documentation.

## get-playstore-apps-by-pname-sdks

`GET /playstore/apps/{pname}/sdks` — Play Store Apps » SDKs

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `pname` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=500 |
| `sort_field` | query | false | default=version_code_min |
| `sort_order` | query | false | default=desc |

Response fields: `name` (str), `slug` (str), `upload_date_min` (date), `upload_date_max` (date), `version_code_min` (int), `version_code_max` (int), `currently_installed` (bool).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/apps/{pname}/sdks). Verified 2026-09-09 signed-in documentation.

## get-playstore-apps-by-pname-versions

`GET /playstore/apps/{pname}/versions` — Play Store Apps » Versions

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `pname` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=version_code |
| `sort_order` | query | false | default=desc |

Response fields: `id` (int), `asset_size` (int), `version_string` (str), `downloaded` (bool), `is_latest` (bool), `latest_since` (datetime), `release_date` (date), `release_notes` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/apps/{pname}/versions). Verified 2026-09-09 signed-in documentation.

## get-playstore-apps-by-pname-versions-by-version-code-apks

`GET /playstore/apps/{pname}/versions/{version_code}/apks` — Play Store Apps » APKs

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `pname` | path | true |  |
| `version_code` | path | true |  |

Response fields: `id` (int), `name` (str), `file_extension` (str), `byte_size` (int), `download_at` (date), `md5` (str), `sha1` (str), `sha256` (str), `namespace_count` (int), `permission_count` (int), `s3uri` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/apps/{pname}/versions/{version_code}/apks). Verified 2026-09-09 signed-in documentation.

## get-playstore-developers-by-developer-id-apps

`GET /playstore/developers/{developer_id}/apps` — Play Store Developers » Apps

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `developer_id` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=downloads |
| `sort_order` | query | false | default=desc |

Response fields: `pname` (str), `title` (str), `rank` (int), `score` (float), `developer_id` (int), `developer_name` (str), `developer_email` (str), `developer_website` (str), `downloads` (int), `removed` (bool), `review_count` (int), `rating` (float), `rating_count` (int), `upload_date` (date), `release_date` (date), `small_icon_url` (str), `large_icon_url` (str), `category` (str), `categories` (str), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `version` (object), `version.major_version_number` (int), `version.version_code` (int), `version.version_string` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/developers/{developer_id}/apps). Verified 2026-09-09 signed-in documentation.

## get-playstore-rankings

`GET /playstore/rankings` — Play Store Rankings » Lists

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



Response fields: `container` (str), `category` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/rankings). Verified 2026-09-09 signed-in documentation.

## get-playstore-rankings-by-container-by-category-by-country

`GET /playstore/rankings/{container}/{category}/{country}` — Play Store Rankings » Rankings

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=ranking |
| `sort_order` | query | false | default=asc |
| `container` | path | true |  |
| `category` | path | true |  |
| `country` | path | true |  |
| `before` | query | false |  |

Response fields: `pname` (str), `title` (str), `rank` (int), `score` (float), `developer_id` (int), `developer_name` (str), `developer_email` (str), `developer_website` (str), `downloads` (int), `removed` (bool), `review_count` (int), `rating` (float), `rating_count` (int), `upload_date` (date), `release_date` (date), `small_icon_url` (str), `large_icon_url` (str), `category` (str), `categories` (str), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `version` (object), `version.major_version_number` (int), `version.version_code` (int), `version.version_string` (str), `ranking` (int), `ranking_date` (datetime).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/rankings/{container}/{category}/{country}). Verified 2026-09-09 signed-in documentation.

## get-playstore-sdks

`GET /playstore/sdks` — Play Store SDKs » Directory

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `offset` | query | false | default=0; min=0 |
| `company.id` | query | false |  |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=rank |
| `sort_order` | query | false | default=desc |

Response fields: `rank` (int), `score` (float), `name` (str), `id` (int), `slug` (str), `total_installs` (int), `current_installs` (int), `uninstalls` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/sdks). Verified 2026-09-09 signed-in documentation.

## get-playstore-sdks-by-sdk-slug

`GET /playstore/sdks/{sdk_slug}` — Play Store SDKs » Overview

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `sdk_slug` | path | true |  |

Response fields: `name` (str), `slug` (str), `description` (str), `tags` (list[str]), `url` (str), `logo_url` (str), `publisher_count` (int), `total_installs` (int), `current_installs` (int), `uninstalls` (int), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/sdks/{sdk_slug}). Verified 2026-09-09 signed-in documentation.

## get-playstore-sdks-by-sdk-slug-categories

`GET /playstore/sdks/{sdk_slug}/categories` — Play Store SDKs » Categories

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `sdk_slug` | path | true |  |

Response fields: `category` (str), `total_installs` (int), `current_installs` (int), `uninstalls` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/sdks/{sdk_slug}/categories). Verified 2026-09-09 signed-in documentation.

## get-playstore-sdks-by-sdk-slug-downloads

`GET /playstore/sdks/{sdk_slug}/downloads` — Play Store SDKs » Downloads

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `sdk_slug` | path | true |  |

Response fields: `downloads` (int), `total_installs` (int), `current_installs` (int), `uninstalls` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/sdks/{sdk_slug}/downloads). Verified 2026-09-09 signed-in documentation.

## get-playstore-sdks-by-sdk-slug-install-trend

`GET /playstore/sdks/{sdk_slug}/install-trend` — Play Store SDKs » Install trend

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `sdk_slug` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |

Response fields: `total_installs` (int), `current_installs` (int), `uninstalls` (int), `date` (datetime).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/sdks/{sdk_slug}/install-trend). Verified 2026-09-09 signed-in documentation.

## get-playstore-sdks-by-sdk-slug-installs

`GET /playstore/sdks/{sdk_slug}/installs` — Play Store SDKs » Installs

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `sdk_slug` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=downloads |
| `sort_order` | query | false | default=desc |

Response fields: `pname` (str), `title` (str), `rank` (int), `score` (float), `developer_id` (int), `developer_name` (str), `developer_email` (str), `developer_website` (str), `downloads` (int), `removed` (bool), `review_count` (int), `rating` (float), `rating_count` (int), `upload_date` (date), `release_date` (date), `small_icon_url` (str), `large_icon_url` (str), `category` (str), `categories` (str), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `version` (object), `version.major_version_number` (int), `version.version_code` (int), `version.version_string` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/sdks/{sdk_slug}/installs). Verified 2026-09-09 signed-in documentation.

## get-playstore-sdks-by-sdk-slug-uninstalls

`GET /playstore/sdks/{sdk_slug}/uninstalls` — Play Store SDKs » Uninstalls

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `sdk_slug` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=downloads |
| `sort_order` | query | false | default=desc |

Response fields: `pname` (str), `title` (str), `rank` (int), `score` (float), `developer_id` (int), `developer_name` (str), `developer_email` (str), `developer_website` (str), `downloads` (int), `removed` (bool), `review_count` (int), `rating` (float), `rating_count` (int), `upload_date` (date), `release_date` (date), `small_icon_url` (str), `large_icon_url` (str), `category` (str), `categories` (str), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `version` (object), `version.major_version_number` (int), `version.version_code` (int), `version.version_string` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/playstore/sdks/{sdk_slug}/uninstalls). Verified 2026-09-09 signed-in documentation.

## get-segment-features-by-segment

`GET /segment-features/{segment}` — Audience Segments » Features

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `segment` | path | true |  |
| `type` | query | true | enum=['Persons', 'Companies'] |
| `limit` | query | false | default=100; min=1 |
| `search` | query | false |  |
| `code` | query | false |  |
| `company_id` | query | false |  |
| `url` | query | false |  |
| `org_id` | query | false |  |
| `linkedin_url` | query | false |  |
| `expand_search` | query | false |  |
| `concat_features` | query | false |  |
| `tag_type` | query | false |  |
| `modifiers` | query | false | enum=['Persons', 'Companies'] |

Response fields: `id` (int), `mod` (str), `name` (str), `description` (str), `explanation` (str), `feature_concat` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/segment-features/{segment}). Verified 2026-09-09 signed-in documentation.

## get-segment-result

`GET /segment-result` — Audience Segments » Segment Result

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).

200 returns a Zstandard-compressed file; 204 means pending, 403 denied, 404 missing. Preserve the token and poll the result, not the submission.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `token` | query | true |  |

Response fields: `status` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/segment-result). Verified 2026-09-09 signed-in documentation.

## get-twitter-query-by-query-ids

`GET /twitter/query/{query_ids}` — Twitter » Twitter Query Status

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `query_ids` | path | true |  |

Response fields: `id` (int), `page_limit` (int), `request_count` (int), `first_seen` (int), `last_seen` (int), `metadata` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/twitter/query/{query_ids}). Verified 2026-09-09 signed-in documentation.

## get-usage

`GET /usage` — Account » Usage

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `month` | query | false |  |
| `impersonate` | query | false |  |

Response fields: `account_name` (str), `total_events` (int), `total_data_quantity` (int), `details` (list[object]), `details[].events` (int), `details[].data_quantity` (int), `details[].category` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/usage). Verified 2026-09-09 signed-in documentation.

## get-web-sites

`GET /web/sites` — Websites » Directory

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=rank |
| `company.id` | query | false |  |
| `search` | query | false |  |

Response fields: `domain` (str), `rank` (int), `alexa_rank` (int), `score` (float), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/web/sites). Verified 2026-09-09 signed-in documentation.

## get-web-sites-by-domain

`GET /web/sites/{domain}` — Websites » Overview

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `domain` | path | true |  |

Response fields: `domain` (str), `rank` (int), `alexa_rank` (int), `score` (float), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/web/sites/{domain}). Verified 2026-09-09 signed-in documentation.

## get-web-sites-by-domain-tags

`GET /web/sites/{domain}/tags` — Websites » Tags

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `domain` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |

Response fields: `id` (int), `slug` (str), `name` (str), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `times_seen` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/web/sites/{domain}/tags). Verified 2026-09-09 signed-in documentation.

## get-web-tags

`GET /web/tags` — Web Tags » Directory

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=rank |
| `sort_order` | query | false | default=desc |

Response fields: `rank` (int), `score` (float), `name` (str), `id` (int), `slug` (str), `domain_count` (int), `times_seen` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/web/tags). Verified 2026-09-09 signed-in documentation.

## get-web-tags-by-slug

`GET /web/tags/{slug}` — Web Tags » Overview

Encoding: query. Pagination: none. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `slug` | path | true |  |

Response fields: `id` (int), `slug` (str), `name` (str), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str), `description` (str), `url` (str), `logo_url` (str), `site_count` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/web/tags/{slug}). Verified 2026-09-09 signed-in documentation.

## get-web-tags-by-slug-sites

`GET /web/tags/{slug}/sites` — Web Tags » Sites

Encoding: query. Pagination: offset. Automatic replay eligible: true (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `slug` | path | true |  |
| `offset` | query | false | default=0; min=0 |
| `page_size` | query | false | default=4; min=1; max=100 |
| `sort_field` | query | false | default=rank |

Response fields: `domain` (str), `rank` (int), `alexa_rank` (int), `score` (float), `company` (object), `company.id` (int), `company.name` (str), `company.slug` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/web/tags/{slug}/sites). Verified 2026-09-09 signed-in documentation.

## post-companies-segment-preview

`POST /companies/segment-preview` — Companies » Company Segment Preview

Encoding: form. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `query` | form | true |  |
| `limit` | form | false | default=100; min=1; max=100 |

Response fields: `total` (int), `is_truncated` (int), `results` (list[int]).

Sources: [provider contract](https://mixrank.com/api/documentation#/companies/segment-preview). Verified 2026-09-09 signed-in documentation.

## post-companies-segment-query

`POST /companies/segment-query` — Companies » Company Segment Query

Encoding: form. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `query` | form | true |  |

Response fields: `token` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/companies/segment-query). Verified 2026-09-09 signed-in documentation.

## post-email-validate-bulk-job

`POST /email/validate/bulk-job` — Email » Submit Job

Encoding: multipart. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).

POST uploads UTF-8 newline-separated emails as multipart file; limit 3000000 lines and 20 active jobs. private=true requires strategy=fetch and explicit deletion of downloaded output.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `maxage` | query | false |  |
| `strategy` | query | false | default=besteffort; enum=['cached', 'fetch', 'strict', 'besteffort'] |
| `private` | query | false | default=False |
| `file` | file | true |  |

Response fields: `id` (str), `created_at` (datetime), `lines_total` (int), `lines_invalid` (int), `lines_processed` (int), `lines_refreshed` (int), `strategy` (str), `maxage` (datetime), `completed_at` (datetime), `download_url` (str), `private` (bool).

Sources: [provider contract](https://mixrank.com/api/documentation#/email/validate/bulk-job). Verified 2026-09-09 signed-in documentation.

## post-identity-ingest

`POST /identity/ingest` — Identity » Ingest

Encoding: form. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).

Contributes records to the provider identity graph. accounts is JSON encoded in a form field; source and source_key identify the contribution.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `accounts` | form | true | enum=['email_address', 'phone_number', 'linkedin_profile', 'linkedin_company', 'domain', 'soundcloud_profile', 'klout_profile', 'goodreads_profile', 'foursquare_profile', 'flickr_profile', 'vimeo_profile', 'facebook_profile', 'spotify_artist', 'blogger_profile', 'meetup_profile', 'stackoverflow_profile', 'itunes_artist', 'gplus_profile', 'twitter_profile', 'pinterest_profile', 'github_profile', 'tripit_profile', 'aboutme_profile', 'angellist_profile', 'behance_profile', 'crunchbase_profile', 'dribbble_profile', 'instagram_profile', 'myspace_profile', 'quora_profile', 'reddit_profile', 'snapchat_profile', 'tiktok_profile', 'tumblr_profile', 'xing_profile', 'play_store_developer', 'linktree_profile', 'crunchbase_company', 'clubhouse_profile', 'xing_company', 'gravatar_profile', 'whatsapp_profile', 'skype_profile', 'telegram_profile', 'play_app', 'ios_app'] |
| `source` | form | true |  |
| `source_key` | form | true |  |

Sources: [provider contract](https://mixrank.com/api/documentation#/identity/ingest). Verified 2026-09-09 signed-in documentation.

## post-linkedin-post-bulk-job

`POST /linkedin/post/bulk_job` — LinkedIn » Post Bulk Scan - Submit Job

Encoding: multipart. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).

POST uploads UTF-8 post URLs as multipart file; limit 10000 lines and 20 active jobs. Submission cannot be cancelled. Use returned job id to inspect completion.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `file` | file | true |  |

Response fields: `id` (str), `created_at` (datetime), `lines_total` (int), `lines_invalid` (int), `lines_processed` (int), `completed_at` (datetime).

Sources: [provider contract](https://mixrank.com/api/documentation#/linkedin/post/bulk_job). Verified 2026-09-09 signed-in documentation.

## post-linkedin-profile-bulk-job

`POST /linkedin/profile/bulk_job` — LinkedIn » Submit Job

Encoding: multipart. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).

POST uploads UTF-8 profile URLs as multipart file; limit 250000 lines and 20 active jobs. Submitted jobs cannot be cancelled. completed_at and download_url indicate completion. Toolkit defaults to cached; fetch/strict/besteffort need explicit refresh. Synchronous LiveScan returns profile data; 512 is temporary fetch failure, 513 is permanent. No implicit strategy fallback.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `maxage` | query | false |  |
| `strategy` | query | false | default=fetch; enum=['fetch', 'cached', 'strict', 'besteffort'] |
| `file` | file | true |  |

Response fields: `id` (str), `created_at` (datetime), `lines_total` (int), `lines_invalid` (int), `lines_processed` (int), `lines_refreshed` (int), `lines_cached` (int), `maxage` (int), `strategy` (str), `completed_at` (datetime), `download_url` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/linkedin/profile/bulk_job). Verified 2026-09-09 signed-in documentation.

## post-person-segment-preview

`POST /person/segment-preview` — People » Person Segment Preview

Encoding: form. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `query` | form | true |  |
| `limit` | form | false | default=100; min=1; max=100 |

Response fields: `total` (int), `is_truncated` (int), `results` (list[int]).

Sources: [provider contract](https://mixrank.com/api/documentation#/person/segment-preview). Verified 2026-09-09 signed-in documentation.

## post-person-segment-query

`POST /person/segment-query` — People » Person Segment Query

Encoding: form. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `query` | form | true |  |

Response fields: `token` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/person/segment-query). Verified 2026-09-09 signed-in documentation.

## post-privacy-redaction-request

`POST /privacy-redaction-request` — Privacy » Privacy Redact

Encoding: form. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).

Submits a privacy redaction request; external mutation requiring explicit task authorization.

| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `linkedin_url` | form | false |  |

Sources: [provider contract](https://mixrank.com/api/documentation#/privacy-redaction-request). Verified 2026-09-09 signed-in documentation.

## post-twitter-profile-by-handler

`POST /twitter/profile/{handler}` — Twitter » Twitter Profile

Encoding: form. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `handler` | path | true |  |
| `interval` | form | false |  |

Response fields: `username` (int), `interval` (str).

Sources: [provider contract](https://mixrank.com/api/documentation#/twitter/profile/{handler}). Verified 2026-09-09 signed-in documentation.

## post-twitter-query-by-query

`POST /twitter/query/{query}` — Twitter » Twitter Query

Encoding: form. Pagination: none. Automatic replay eligible: false (refresh can further restrict this).



| Parameter | Location | Required | Default / range / values |
|---|---|---|---|
| `query` | path | true |  |
| `search_key` | form | true |  |
| `search_type` | form | false | default=1 |
| `page_limit` | form | false |  |
| `disabled` | form | false |  |
| `interval` | form | false |  |

Response fields: `id` (int), `page_limit` (int), `request_count` (int), `first_seen` (int), `last_seen` (int), `metadata` (int).

Sources: [provider contract](https://mixrank.com/api/documentation#/twitter/query/{query}). Verified 2026-09-09 signed-in documentation.

## search-companies

`POST /elasticsearch/companies/_search` — Search companies with Elasticsearch DSL

Encoding: json. Pagination: elasticsearch. Automatic replay eligible: true (refresh can further restrict this).

Pass the complete DSL body, including nested queries, _source, sort, aggregations and autocomplete. Inspect bundled field mappings before constructing queries. No server-version assumption or undocumented mapping endpoint.

Response fields: `hits` (object), `aggregations` (object).

Sources: [provider contract](https://docs.google.com/document/d/1PyTnOfF__OE4HedeIl2L_XzhjhbG_VsC8M8H5mWdKQ0/edit). Verified 2026-09-09 provider document.

## search-industries

`POST /elasticsearch/industries/_search` — Search industries with Elasticsearch DSL

Encoding: json. Pagination: elasticsearch. Automatic replay eligible: true (refresh can further restrict this).

Pass the complete DSL body, including nested queries, _source, sort, aggregations and autocomplete. Inspect bundled field mappings before constructing queries. No server-version assumption or undocumented mapping endpoint.

Response fields: `hits` (object), `aggregations` (object).

Sources: [provider contract](https://docs.google.com/document/d/1PyTnOfF__OE4HedeIl2L_XzhjhbG_VsC8M8H5mWdKQ0/edit). Verified 2026-09-09 provider document.

## search-jobs

`POST /elasticsearch/jobs/_search` — Search jobs with Elasticsearch DSL

Encoding: json. Pagination: elasticsearch. Automatic replay eligible: true (refresh can further restrict this).

Pass the complete DSL body, including nested queries, _source, sort, aggregations and autocomplete. Inspect bundled field mappings before constructing queries. No server-version assumption or undocumented mapping endpoint.

Response fields: `hits` (object), `aggregations` (object).

Sources: [provider contract](https://docs.google.com/document/d/1PyTnOfF__OE4HedeIl2L_XzhjhbG_VsC8M8H5mWdKQ0/edit). Verified 2026-09-09 provider document.

## search-org-name

`POST /elasticsearch/org_name/_search` — Search org_name with Elasticsearch DSL

Encoding: json. Pagination: elasticsearch. Automatic replay eligible: true (refresh can further restrict this).

Pass the complete DSL body, including nested queries, _source, sort, aggregations and autocomplete. Inspect bundled field mappings before constructing queries. No server-version assumption or undocumented mapping endpoint.

Response fields: `hits` (object), `aggregations` (object).

Sources: [provider contract](https://docs.google.com/document/d/1PyTnOfF__OE4HedeIl2L_XzhjhbG_VsC8M8H5mWdKQ0/edit). Verified 2026-09-09 provider document.

## search-person2

`POST /elasticsearch/person2/_search` — Search person2 with Elasticsearch DSL

Encoding: json. Pagination: elasticsearch. Automatic replay eligible: true (refresh can further restrict this).

Pass the complete DSL body, including nested queries, _source, sort, aggregations and autocomplete. Inspect bundled field mappings before constructing queries. No server-version assumption or undocumented mapping endpoint.

Response fields: `hits` (object), `aggregations` (object).

Sources: [provider contract](https://docs.google.com/document/d/1PyTnOfF__OE4HedeIl2L_XzhjhbG_VsC8M8H5mWdKQ0/edit). Verified 2026-09-09 provider document.

