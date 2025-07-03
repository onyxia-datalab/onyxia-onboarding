# Changelog

## [0.1.7](https://github.com/onyxia-datalab/onyxia-onboarding/compare/v0.1.6...v0.1.7) (2025-07-03)


### Features

* add contextPath ([#65](https://github.com/onyxia-datalab/onyxia-onboarding/issues/65)) ([2359a91](https://github.com/onyxia-datalab/onyxia-onboarding/commit/2359a91c2031a435bc003479decd1bb4213fd93b))


### Bug Fixes

* **deps:** update go minor and patch updates ([#63](https://github.com/onyxia-datalab/onyxia-onboarding/issues/63)) ([add1a05](https://github.com/onyxia-datalab/onyxia-onboarding/commit/add1a05661e3024481f6c56e531aee438568b1b3))
* **deps:** update module github.com/go-chi/chi/v5 to v5.2.2 [security] ([#64](https://github.com/onyxia-datalab/onyxia-onboarding/issues/64)) ([ef95327](https://github.com/onyxia-datalab/onyxia-onboarding/commit/ef95327a6cba6102d44d930f8f4b8df382d4d997))

## [0.1.6](https://github.com/onyxia-datalab/onyxia-onboarding/compare/v0.1.5...v0.1.6) (2025-06-19)


### Features

* add helm chart and oci push action ([#48](https://github.com/onyxia-datalab/onyxia-onboarding/issues/48)) ([40547f0](https://github.com/onyxia-datalab/onyxia-onboarding/commit/40547f04f8125991ef3865529e4e15d7890b383e))


### Bug Fixes

* **Dockerfile:** app path & go compile option for distroless image ([#61](https://github.com/onyxia-datalab/onyxia-onboarding/issues/61)) ([3de4f94](https://github.com/onyxia-datalab/onyxia-onboarding/commit/3de4f945035fd8692c93b12e2196fd41d7a10c25))

## [0.1.5](https://github.com/onyxia-datalab/onyxia-onboarding/compare/v0.1.4...v0.1.5) (2025-02-27)


### Features

* trigger release ([07ecbac](https://github.com/onyxia-datalab/onyxia-onboarding/commit/07ecbac285eb029c2e64d36946903a746d4faa77))

## [0.1.4](https://github.com/onyxia-datalab/onyxia-onboarding/compare/v0.1.3...v0.1.4) (2025-02-27)


### Features

* trigger release ([17c9b7e](https://github.com/onyxia-datalab/onyxia-onboarding/commit/17c9b7e6dde1a184bcf62fc86be3668f6e01ccf4))

## [0.1.3](https://github.com/onyxia-datalab/onyxia-onboarding/compare/v0.1.2...v0.1.3) (2025-02-27)


### Features

* add qemu and buildx setup to reduce docker actions time ([#44](https://github.com/onyxia-datalab/onyxia-onboarding/issues/44)) ([e56c3b6](https://github.com/onyxia-datalab/onyxia-onboarding/commit/e56c3b63e32193d9256b329396a731c3eb94cc4d))

## [0.1.2](https://github.com/onyxia-datalab/onyxia-onboarding/compare/v0.1.1...v0.1.2) (2025-02-27)


### Features

* trigger release ([ad028b6](https://github.com/onyxia-datalab/onyxia-onboarding/commit/ad028b618ff25dc1b0dda5649a2e0cca17609691))

## [0.1.1](https://github.com/onyxia-datalab/onyxia-onboarding/compare/v0.1.0...v0.1.1) (2025-02-25)


### Bug Fixes

* **ci:** make tags available in ci so docker tags are correct ([#38](https://github.com/onyxia-datalab/onyxia-onboarding/issues/38)) ([dfaa2dc](https://github.com/onyxia-datalab/onyxia-onboarding/commit/dfaa2dc9cbd85668da5944ba506dcf50588e0949))

## 0.1.0 (2025-02-25)


### Features

* add chi first route ([#8](https://github.com/onyxia-datalab/onyxia-onboarding/issues/8)) ([6e2af0a](https://github.com/onyxia-datalab/onyxia-onboarding/commit/6e2af0ad987a564890880b42bb0b6f076d3802f8))
* add quotas support for namespace, split files  ([#25](https://github.com/onyxia-datalab/onyxia-onboarding/issues/25)) ([0fa9e89](https://github.com/onyxia-datalab/onyxia-onboarding/commit/0fa9e899738c5bf04d891132a16e50fbec09ded6))
* add support of env variables ([#19](https://github.com/onyxia-datalab/onyxia-onboarding/issues/19)) ([37ffbd4](https://github.com/onyxia-datalab/onyxia-onboarding/commit/37ffbd4469e0f102bd9f92efed69fcb9df0425ef))
* add username, groups and roles in log ([2b0b5cc](https://github.com/onyxia-datalab/onyxia-onboarding/commit/2b0b5cc2f76a1d819bdf81b665a25b6f366d3521))
* clean archi with ogen ([#21](https://github.com/onyxia-datalab/onyxia-onboarding/issues/21)) ([a1cb014](https://github.com/onyxia-datalab/onyxia-onboarding/commit/a1cb0140b922bb767405a409a8b48fde38795221))
* Implement role-based quotas and validate group onboarding rights ([#33](https://github.com/onyxia-datalab/onyxia-onboarding/issues/33)) ([d61ad17](https://github.com/onyxia-datalab/onyxia-onboarding/commit/d61ad171cc9e96af007554e4be9ce8efb8eb81d5))
* Improve default env handling with embedded config ([#34](https://github.com/onyxia-datalab/onyxia-onboarding/issues/34)) ([ddc79b2](https://github.com/onyxia-datalab/onyxia-onboarding/commit/ddc79b22025af30969aeef1c3b0da1cd7ae4a0e8))
* makefile and adapt CI ([#36](https://github.com/onyxia-datalab/onyxia-onboarding/issues/36)) ([4cdfcc4](https://github.com/onyxia-datalab/onyxia-onboarding/commit/4cdfcc4e9d3984b7e9a04691f5c7887c4eaaacba))
* role base quotas for user and refacto ctx ([#35](https://github.com/onyxia-datalab/onyxia-onboarding/issues/35)) ([b5bca29](https://github.com/onyxia-datalab/onyxia-onboarding/commit/b5bca29ddbf3be27d64cd04dcd4211a661b4256a))
* setup renovate ([#4](https://github.com/onyxia-datalab/onyxia-onboarding/issues/4)) ([96859c4](https://github.com/onyxia-datalab/onyxia-onboarding/commit/96859c441696bd88745ba420fb20a0f9770621f6))


### Bug Fixes

* **deps:** update kubernetes packages to v0.32.2 ([#29](https://github.com/onyxia-datalab/onyxia-onboarding/issues/29)) ([5c6a47f](https://github.com/onyxia-datalab/onyxia-onboarding/commit/5c6a47fba4a9689ee863216ac77cd6d7594fc2ad))
* **deps:** update module github.com/ogen-go/ogen to v1.10.0 ([#24](https://github.com/onyxia-datalab/onyxia-onboarding/issues/24)) ([963aaef](https://github.com/onyxia-datalab/onyxia-onboarding/commit/963aaef99ad611c33f1e77017491f2b58131019f))
* error introduced by [#8](https://github.com/onyxia-datalab/onyxia-onboarding/issues/8) ([cb53a31](https://github.com/onyxia-datalab/onyxia-onboarding/commit/cb53a310dc53ecaf22cdf3986349c39fd7ebd677))
* humbly fixing the linting error ([#13](https://github.com/onyxia-datalab/onyxia-onboarding/issues/13)) ([f9b9d24](https://github.com/onyxia-datalab/onyxia-onboarding/commit/f9b9d2409397d76b83d552f989b8f1ebbb3420aa))
* oidc groups and roles extractions ([ab72ac2](https://github.com/onyxia-datalab/onyxia-onboarding/commit/ab72ac297bd44aa68e79939d89de760879b83de1))
* oidc_test ([7c0f35e](https://github.com/onyxia-datalab/onyxia-onboarding/commit/7c0f35ee03b9485e34fff5c1e2d670b27f1c8d44))
* renovate use conventional commits ([#7](https://github.com/onyxia-datalab/onyxia-onboarding/issues/7)) ([456e7b1](https://github.com/onyxia-datalab/onyxia-onboarding/commit/456e7b112aaa7e37b0785c96847780cc43406e05))
* **test:** ignore renovate PRs ([69083bc](https://github.com/onyxia-datalab/onyxia-onboarding/commit/69083bc6048b96b58cea2d06af0185698a1add1a))
