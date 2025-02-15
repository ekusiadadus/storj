// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

<template>
    <div>
        <div class="access-grant__modal-container__header-container">
            <h2 class="access-grant__modal-container__header-container__title">Create Access</h2>
            <div
                class="access-grant__modal-container__header-container__close-cross-container" @click="onCloseClick"
            >
                <CloseCrossIcon />
            </div>
        </div>
        <div class="access-grant__modal-container__body-container">
            <TypesIcon class="access-grant__modal-container__body-container__type-icon" />
            <div class="access-grant__modal-container__body-container__type">
                <p>Type</p>
                <div class="access-grant__modal-container__body-container__type__type-container">
                    <input
                        id="access-grant-check"
                        v-model="checkedType"
                        value="access"
                        type="radio"
                        :checked="checkedType === 'access'"
                        name="type"
                    >
                    <label for="access-grant-check">
                        Access Grant
                    </label>
                    <img
                        class="tooltip-icon"
                        src="../../../../static/images/accessGrants/create-access_information.png"
                        alt="tooltip icon"
                        @mouseover="toggleTooltipHover('access','over')"
                        @mouseleave="toggleTooltipHover('access','leave')"
                    >
                    <div
                        v-if="tooltipHover === 'access'"
                        class="access-tooltip"
                        @mouseover="toggleTooltipHover('access','over')"
                        @mouseleave="toggleTooltipHover('access','leave')"
                    >
                        <span class="tooltip-text">Keys to upload, delete, and view your project's data.  <a class="tooltip-link" href="https://storj-labs.gitbook.io/dcs/concepts/access/access-grants" target="_blank" rel="noreferrer noopener">Learn More</a></span>
                    </div>
                </div>
                <div class="access-grant__modal-container__body-container__type__type-container">
                    <input
                        id="s3-check"
                        v-model="checkedType"
                        value="s3"
                        type="radio"
                        name="type"
                        :checked="checkedType === 's3'"
                    >
                    <label for="s3-check">S3 Credentials</label>
                    <img
                        class="tooltip-icon"
                        src="../../../../static/images/accessGrants/create-access_information.png"
                        alt="tooltip icon"
                        @mouseover="toggleTooltipHover('s3','over')"
                        @mouseleave="toggleTooltipHover('s3','leave')"
                    >
                    <div
                        v-if="tooltipHover === 's3'"
                        class="s3-tooltip"
                        @mouseover="toggleTooltipHover('s3','over')"
                        @mouseleave="toggleTooltipHover('s3','leave')"
                    >
                        <span class="tooltip-text">Generates access key, secret key, and endpoint to use in your S3-supporting application.  <a class="tooltip-link" href="https://docs.storj.io/dcs/api-reference/s3-compatible-gateway" target="_blank" rel="noreferrer noopener">Learn More</a></span>
                    </div>
                </div>
                <div class="access-grant__modal-container__body-container__type__type-container">
                    <input
                        id="api-check"
                        v-model="checkedType"
                        value="api"
                        type="radio"
                        name="type"
                        :checked="checkedType === 'api'"
                    >
                    <label for="api-check">API Access</label>
                    <img
                        class="tooltip-icon"
                        src="../../../../static/images/accessGrants/create-access_information.png"
                        alt="tooltip icon"
                        @mouseover="toggleTooltipHover('api','over')"
                        @mouseleave="toggleTooltipHover('api','leave')"
                    >
                    <div
                        v-if="tooltipHover === 'api'"
                        class="api-tooltip"
                        @mouseover="toggleTooltipHover('api','over')"
                        @mouseleave="toggleTooltipHover('api','leave')"
                    >
                        <span class="tooltip-text">Creates access grant to run in the command line.  <a class="tooltip-link" href="https://docs.storj.io/dcs/getting-started/quickstart-uplink-cli/generate-access-grants-and-tokens/generate-a-token/" target="_blank" rel="noreferrer noopener">Learn More</a></span>
                    </div>
                </div>
            </div>
            <NameIcon class="access-grant__modal-container__body-container__name-icon" />
            <div class="access-grant__modal-container__body-container__name">
                <p>Name</p>
                <input
                    v-model="accessName"
                    type="text"
                    placeholder="Input Access Name" class="access-grant__modal-container__body-container__name__input"
                >
            </div>
            <PermissionsIcon class="access-grant__modal-container__body-container__permissions-icon" />
            <div class="access-grant__modal-container__body-container__permissions">
                <p>Permissions</p>
                <div>
                    <input
                        id="permissions__all-check"
                        type="checkbox"
                        :checked="allPermissionsClicked"
                        @click="toggleAllPermission('all')"
                    >
                    <label for="permissions__all-check">All</label>
                    <Chevron :class="`permissions-chevron-${showAllPermissions.position}`" @click="togglePermissions" />
                </div>
                <div v-if="showAllPermissions.show === true">
                    <div v-for="(item, key) in permissionsList" :key="key">
                        <input
                            :id="`permissions__${item}-check`"
                            v-model="selectedPermissions"
                            :value="item"
                            type="checkbox"
                            :checked="checkedPermissions.item"
                            @click="toggleAllPermission(item)"
                        >
                        <label :for="`permissions__${item}-check`">{{ item }}</label>
                    </div>
                </div>
            </div>
            <BucketsIcon class="access-grant__modal-container__body-container__buckets-icon" />
            <div class="access-grant__modal-container__body-container__buckets">
                <p>Buckets</p>
                <div>
                    <BucketsSelection
                        class="access-bucket-container"
                        :show-scrollbar="true"
                    />
                </div>
                <div class="access-grant__modal-container__body-container__buckets__bucket-bullets">
                    <div
                        v-for="(name, index) in selectedBucketNames"
                        :key="index"
                        class="access-grant__modal-container__body-container__buckets__bucket-bullets__container"
                    >
                        <BucketNameBullet :name="name" />
                    </div>
                </div>
            </div>
            <DateIcon class="access-grant__modal-container__body-container__date-icon" />
            <div class="access-grant__modal-container__body-container__duration">
                <p>Duration</p>
                <div v-if="addDateSelected">
                    <DurationSelection
                        container-style="access-date-container"
                        text-style="access-date-text"
                        picker-style="__access-date-container"
                    />
                </div>
                <div
                    v-else
                    class="access-grant__modal-container__body-container__duration__text"
                    @click="addDateSelected = true"
                >
                    Add Date (optional)
                </div>
            </div>

        <!-- for future use when notes feature is implemented -->
        <!-- <NotesIcon class="access-grant__modal-container__body-container__notes-icon"/>
                    <div class="access-grant__modal-container__body-container__notes">
                        <p>Notes</p>
                        <div>--Notes Section Here--</div>
                    </div> -->
        </div>
        <div class="access-grant__modal-container__footer-container">
            <a href="https://docs.storj.io/dcs/concepts/access/access-grants/api-key/" target="_blank" rel="noopener noreferrer">
                <v-button
                    label="Learn More"
                    width="150px"
                    height="50px"
                    is-transparent="true"
                    font-size="16px"
                    class="access-grant__modal-container__footer-container__learn-more-button"
                />
            </a>
            <v-button
                :label="checkedType === 'api' ? 'Create Keys  ⟶' : 'Encrypt My Access  ⟶'"
                font-size="16px"
                width="auto"
                height="50px"
                class="access-grant__modal-container__footer-container__encrypt-button"
                :on-press="checkedType === 'api' ? propogateInfo : encryptClickAction"
                :is-disabled="selectedPermissions.length === 0 || accessName === ''"
            />
        </div>
    </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import DateIcon from '@/../static/images/accessGrants/create-access_date.svg';
import VButton from '@/components/common/VButton.vue';
import BucketsSelection from '@/components/accessGrants/permissions/BucketsSelection.vue';
import CloseCrossIcon from '@/../static/images/common/closeCross.svg';
import TypesIcon from '@/../static/images/accessGrants/create-access_type.svg';
import NameIcon from '@/../static/images/accessGrants/create-access_name.svg';
import PermissionsIcon from '@/../static/images/accessGrants/create-access_permissions.svg';
import Chevron from '@/../static/images/accessGrants/chevron.svg';
import BucketsIcon from '@/../static/images/accessGrants/create-access_buckets.svg';
import BucketNameBullet from "@/components/accessGrants/permissions/BucketNameBullet.vue";
import DurationSelection from '@/components/accessGrants/permissions/DurationSelection.vue';

import { ACCESS_GRANTS_ACTIONS } from '@/store/modules/accessGrants';
import { AccessGrant } from '@/types/accessGrants';


// @vue/component
@Component({
    components: {
        Chevron,
        CloseCrossIcon,
        TypesIcon,
        NameIcon,
        PermissionsIcon,
        BucketsSelection,
        BucketsIcon,
        BucketNameBullet,
        DateIcon,
        DurationSelection,
        VButton
    },
})


export default class CreateFormModal extends Vue {
    @Prop({ default: '' })
    private checkedType: string;

    public showAllPermissions = {show: false, position: "up"};
    private accessName = '';
    private selectedPermissions : string[] = [];
    private allPermissionsClicked = false;
    private permissionsList = ["Read","Write","List","Delete"];
    private checkedPermissions = {Read: false, Write: false, List: false, Delete: false};
    private accessGrantList = this.accessGrantsList;
    private addDateSelected = false;
    public tooltipHover = '';
    public tooltipVisibilityTimer;

    public mounted(): void {
        this.showAllPermissions = {show: false, position: "up"};
    }

    public onCloseClick(): void {
        this.$store.dispatch(ACCESS_GRANTS_ACTIONS.CLEAR_SELECTION);
        this.$emit('close-modal');
    }

    /**
     * Retrieves selected buckets for bucket bullets.
     */
    public get selectedBucketNames(): string[] {
        return this.$store.state.accessGrantsModule.selectedBucketNames;
    }

    /**
     * propogates selected info to parent on flow progression.
     */
    public propogateInfo(): void {
        const payloadObject  = {
            'checkedType': this.checkedType,
            'accessName': this.accessName,
            'selectedPermissions': this.selectedPermissions,
        }

        this.$emit('propogateInfo', payloadObject, this.checkedType)
    }

    /**
     * Toggles permissions list visibility.
     */
    public togglePermissions(): void {
        this.showAllPermissions.show = !this.showAllPermissions.show;
        this.showAllPermissions.position = this.showAllPermissions.show ? 'up' : 'down';
    }

    public get accessGrantsList(): AccessGrant[] {
        return this.$store.state.accessGrantsModule.page.accessGrants;
    }

    public encryptClickAction(): void {
        let mappedList = this.accessGrantList.map((key) => (key.name))
        if (mappedList.includes(this.accessName)) {
            this.$notify.error(`validation: An API Key with this name already exists in this project, please use a different name`);
            return
        } else if (this.checkedType !== "api") {
            // emit event here
            this.propogateInfo()
            this.$emit('encrypt');
        }
    }

    public toggleAllPermission(type): void {
        if (type === 'all' && !this.allPermissionsClicked) {
            this.allPermissionsClicked = true;
            this.selectedPermissions = this.permissionsList;
            this.checkedPermissions = { Read: true, Write: true, List: true, Delete: true }
            return
        } else if(type === 'all' && this.allPermissionsClicked) {
            this.allPermissionsClicked = false;
            this.selectedPermissions = [];
            this.checkedPermissions = { Read: false, Write: false, List: false, Delete: false };
            return
        } else if(this.checkedPermissions[type]) {
            this.checkedPermissions[type] = false;
            this.allPermissionsClicked = false;
            return;
        } else {
            this.checkedPermissions[type] = true;
            if(this.checkedPermissions.Read && this.checkedPermissions.Write && this.checkedPermissions.List && this.checkedPermissions.Delete) {
                this.allPermissionsClicked = true;
            }
        }
    }

    /**
     * Toggles tooltip visibility.
     */
    public toggleTooltipHover(type, action): void {
        if (this.tooltipHover === '' && action === 'over') {
            this.tooltipHover = type;
            return;
        } else if (this.tooltipHover === type && action === 'leave') {
            this.tooltipVisibilityTimer = setTimeout(() => {
                this.tooltipHover = '';
            },750);
            return;
        } else if (this.tooltipHover === type && action === 'over') {
            clearTimeout(this.tooltipVisibilityTimer);
            return;
        } else if(this.tooltipHover !== type) {
            clearTimeout(this.tooltipVisibilityTimer)
            this.tooltipHover = type;
        }
    }

}
</script>

<style scoped lang="scss">
    @mixin chevron {
        padding-left: 4px;
        transition: transform 0.3s;
    }

    @mixin tooltip-container {
        position: absolute;
        background: #56606d;
        border-radius: 6px;
        width: 253px;
        color: #fff;
        display: flex;
        flex-direction: row;
        align-items: flex-start;
        padding: 8px;
        z-index: 1;
        transition: 250ms;
    }

    @mixin tooltip-arrow {
        content: '';
        position: absolute;
        bottom: 0;
        width: 0;
        height: 0;
        border: 6px solid transparent;
        border-top-color: #56606d;
        border-bottom: 0;
        margin-left: -20px;
        margin-bottom: -20px;
    }

    p {
        font-weight: bold;
        padding-bottom: 10px;
    }

    label {
        margin-left: 8px;
        padding-right: 10px;
    }

    .permissions-chevron-up {
        @include chevron;

        transform: rotate(-90deg);
    }

    .permissions-chevron-down {
        @include chevron;
    }

    .tooltip-icon {
        display: flex;
        width: 14px;
        height: 14px;
        cursor: pointer;
    }

    .tooltip-text {
        text-align: center;
        font-weight: 500;
    }

    .access-grant {
        position: fixed;
        top: 0;
        bottom: 0;
        left: 0;
        right: 0;
        z-index: 100;
        background: rgb(27 37 51 / 75%);
        display: flex;
        align-items: flex-start;
        justify-content: center;

        & > * {
            font-family: sans-serif;
        }

        &__modal-container {
            background: #fff;
            border-radius: 6px;
            display: flex;
            flex-direction: column;
            align-items: flex-start;
            position: relative;
            padding: 25px 40px;
            margin-top: 40px;
            width: 410px;
            height: auto;

            &__header-container {
                text-align: left;
                display: grid;
                grid-template-columns: 2fr 1fr;
                width: 100%;
                padding-top: 10px;

                &__title {
                    grid-column: 1;
                }

                &__close-cross-container {
                    grid-column: 2;
                    margin: auto 0 auto auto;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    right: 30px;
                    top: 30px;
                    height: 24px;
                    width: 24px;
                    cursor: pointer;
                }

                &__close-cross-container:hover .close-cross-svg-path {
                    fill: #2683ff;
                }
            }

            &__body-container {
                display: grid;
                grid-template-columns: 1fr 6fr;
                grid-template-rows: auto auto auto auto auto auto;
                grid-row-gap: 24px;
                width: 100%;
                padding-top: 10px;
                margin-top: 24px;

                &__type-icon {
                    grid-column: 1;
                    grid-row: 1;
                }

                &__type {
                    grid-column: 2;
                    grid-row: 1;
                    display: flex;
                    flex-direction: column;

                    &__type-container {
                        display: flex;
                        flex-direction: row;
                        align-items: center;
                        margin-bottom: 10px;
                    }
                }

                &__name-icon {
                    grid-column: 1;
                    grid-row: 2;
                }

                &__name {
                    grid-column: 2;
                    grid-row: 2;
                    display: flex;
                    flex-direction: column;
                    max-width: 238px;

                    &__input {
                        background: #fff;
                        border: 1px solid #c8d3de;
                        box-sizing: border-box;
                        border-radius: 6px;
                        height: 40px;
                        font-size: 17px;
                        padding: 10px;
                    }

                    &__input:focus {
                        border-color: #2683ff;
                    }
                }

                &__input:focus {
                    border-color: #2683ff;
                }

                &__permissions-icon {
                    grid-column: 1;
                    grid-row: 3;
                }

                &__permissions {
                    grid-column: 2;
                    grid-row: 3;
                    display: flex;
                    flex-direction: column;
                }

                &__buckets-icon {
                    grid-column: 1;
                    grid-row: 4;
                }

                &__buckets {
                    grid-column: 2;
                    grid-row: 4;
                    display: flex;
                    flex-direction: column;

                    &__bucket-bullets {
                        display: flex;
                        align-items: center;
                        max-width: 100%;
                        flex-wrap: wrap;

                        &__container {
                            display: flex;
                            margin-top: 5px;
                        }
                    }
                }

                &__date-icon {
                    grid-column: 1;
                    grid-row: 5;
                }

                &__duration {
                    grid-column: 2;
                    grid-row: 5;
                    display: flex;
                    flex-direction: column;

                    &__text {
                        color: #929fb1;
                        text-decoration: underline;
                        font-family: sans-serif;
                        cursor: pointer;
                    }
                }

                &__notes-icon {
                    grid-column: 1;
                    grid-row: 6;
                }

                &__notes {
                    grid-column: 2;
                    grid-row: 6;
                    display: flex;
                    flex-direction: column;
                }
            }

            &__footer-container {
                display: flex;
                width: 100%;
                justify-content: flex-start;
                margin-top: 16px;

                & ::v-deep .container:first-of-type {
                    margin-right: 8px;
                }

                &__learn-more-button {
                    padding: 0 15px;
                }

                &__encrypt-button {
                    padding: 0 15px;
                }
            }
        }
    }

    ::v-deep .buckets-selection {
        margin-left: 0;
        height: 40px;
        border: 1px solid #c8d3de;
    }

    ::v-deep .buckets-selection__toggle-container {
        padding: 10px 20px;
    }

    .access-tooltip {
        top: 68px;
        left: 112px;

        @include tooltip-container;

        &:after {
            left: 50%;
            top: 100%;

            @include tooltip-arrow;
        }
    }

    .s3-tooltip {
        top: 175px;
        left: 121px;

        @include tooltip-container;

        &:after {
            left: 50%;
            top: -8%;
            transform: rotate(180deg);

            @include tooltip-arrow;
        }
    }

    .api-tooltip {
        top: 204px;
        left: 96px;

        @include tooltip-container;

        &:after {
            left: 50%;
            top: -11%;
            transform: rotate(180deg);

            @include tooltip-arrow;
        }
    }

    @media screen and (max-width: 500px) {

        .access-grant__modal-container {
            width: auto;
            max-width: 80vw;
            padding: 30px 24px;

            &__body-container {
                grid-template-columns: 1.2fr 6fr;
            }
        }
    }
</style>