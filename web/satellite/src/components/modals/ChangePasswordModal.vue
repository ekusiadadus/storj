// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

<template>
    <VModal :on-close="closeModal">
        <template #content>
            <div class="change-password">
                <div class="change-password__row">
                    <ChangePasswordIcon />
                    <h2 class="change-password__row__label">Change Password</h2>
                </div>
                <VInput
                    class="full-input"
                    label="Old Password"
                    placeholder="Old Password"
                    is-password="true"
                    :error="oldPasswordError"
                    @setData="setOldPassword"
                />
                <div class="password-input">
                    <VInput
                        class="full-input"
                        label="New Password"
                        placeholder="New Password"
                        is-password="true"
                        :error="newPasswordError"
                        @setData="setNewPassword"
                        @showPasswordStrength="showPasswordStrength"
                        @hidePasswordStrength="hidePasswordStrength"
                    />
                    <PasswordStrength
                        :password-string="newPassword"
                        :is-shown="isPasswordStrengthShown"
                    />
                </div>
                <VInput
                    class="full-input"
                    label="Confirm Password"
                    placeholder="Confirm Password"
                    is-password="true"
                    :error="confirmationPasswordError"
                    @setData="setPasswordConfirmation"
                />
                <div class="change-password__buttons">
                    <VButton
                        label="Cancel"
                        width="100%"
                        height="48px"
                        :on-press="closeModal"
                        is-transparent="true"
                    />
                    <VButton
                        label="Update"
                        width="100%"
                        height="48px"
                        :on-press="onUpdateClick"
                    />
                </div>
            </div>
        </template>
    </VModal>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

import { AuthHttpApi } from '@/api/auth';
import { Validator } from '@/utils/validation';
import { RouteConfig } from "@/router";
import { AnalyticsHttpApi } from '@/api/analytics';
import { AnalyticsEvent } from '@/utils/constants/analyticsEventNames';
import { APP_STATE_MUTATIONS } from "@/store/mutationConstants";

import PasswordStrength from '@/components/common/PasswordStrength.vue';
import VInput from '@/components/common/VInput.vue';
import VButton from '@/components/common/VButton.vue';
import VModal from '@/components/common/VModal.vue';

import ChangePasswordIcon from '@/../static/images/account/changePasswordPopup/changePassword.svg';

const DELAY_BEFORE_REDIRECT = 2000; // 2 sec

// @vue/component
@Component({
    components: {
        ChangePasswordIcon,
        VInput,
        VButton,
        VModal,
        PasswordStrength,
    },
})
export default class ChangePasswordModal extends Vue {
    private oldPassword = '';
    private newPassword = '';
    private confirmationPassword = '';

    private oldPasswordError = '';
    private newPasswordError = '';
    private confirmationPasswordError = '';

    private readonly auth: AuthHttpApi = new AuthHttpApi();

    private readonly analytics: AnalyticsHttpApi = new AnalyticsHttpApi();

    /**
     * Indicates if hint popup needs to be shown while creating new password.
     */
    public isPasswordStrengthShown = false;

    /**
     * Enables password strength info container.
     */
    public showPasswordStrength(): void {
        this.isPasswordStrengthShown = true;
    }

    /**
     * Disables password strength info container.
     */
    public hidePasswordStrength(): void {
        this.isPasswordStrengthShown = false;
    }

    /**
     * Sets old password from input.
     */
    public setOldPassword(value: string): void {
        this.oldPassword = value;
        this.oldPasswordError = '';
    }

    /**
     * Sets new password from input.
     */
    public setNewPassword(value: string): void {
        this.newPassword = value;
        this.newPasswordError = '';
    }

    /**
     * Sets password confirmation from input.
     */
    public setPasswordConfirmation(value: string): void {
        this.confirmationPassword = value;
        this.confirmationPasswordError = '';
    }

    /**
     * Validates inputs and if everything are correct tries to change password and close popup.
     */
    public async onUpdateClick(): Promise<void> {
        let hasError = false;
        if (this.oldPassword.length < 6) {
            this.oldPasswordError = 'Invalid old password. Must be 6 or more characters';
            hasError = true;
        }

        if (!Validator.password(this.newPassword)) {
            this.newPasswordError = 'Invalid password. Use 6 or more characters';
            hasError = true;
        }

        if (!this.confirmationPassword) {
            this.confirmationPasswordError = 'Password required';
            hasError = true;
        }

        if (this.newPassword !== this.confirmationPassword) {
            this.confirmationPasswordError = 'Password not match to new one';
            hasError = true;
        }

        if (hasError) {
            return;
        }

        try {
            await this.auth.changePassword(this.oldPassword, this.newPassword);
        } catch (error) {
            await this.$notify.error(error.message);

            return;
        }

        try {
            await this.auth.logout();

            setTimeout(() => {
                this.$router.push(RouteConfig.Login.path);
            }, DELAY_BEFORE_REDIRECT);
        } catch (error) {
            await this.$notify.error(error.message);
        }

        this.analytics.eventTriggered(AnalyticsEvent.PASSWORD_CHANGED);
        await this.$notify.success('Password successfully changed!');
        this.closeModal();
    }

    /**
     * Closes popup.
     */
    public closeModal(): void {
        this.$store.commit(APP_STATE_MUTATIONS.TOGGLE_CHANGE_PASSWORD_MODAL_SHOWN);
    }
}
</script>

<style scoped lang="scss">
    .change-password {
        background-color: #fff;
        border-radius: 6px;
        display: flex;
        flex-direction: column;
        padding: 48px;

        @media screen and (max-width: 600px) {
            padding: 48px 24px;
        }

        &__row {
            display: flex;
            align-items: center;
            margin-bottom: 20px;

            @media screen and (max-width: 600px) {

                svg {
                    display: none;
                }
            }

            &__label {
                font-family: 'font_bold', sans-serif;
                font-size: 32px;
                line-height: 60px;
                color: #384b65;
                margin: 0 0 0 32px;

                @media screen and (max-width: 600px) {
                    font-size: 24px;
                    line-height: 28px;
                    margin: 0;
                }
            }
        }

        &__buttons {
            width: 100%;
            display: flex;
            align-items: center;
            margin-top: 32px;
            column-gap: 20px;

            @media screen and (max-width: 600px) {
                flex-direction: column-reverse;
                column-gap: unset;
                row-gap: 10px;
                margin-top: 15px;
            }
        }
    }

    .password-input {
        position: relative;
        width: 100%;
    }

    .full-input {
        margin-bottom: 15px;
    }

    @media screen and (max-width: 600px) {

        ::v-deep .password-strength-container {
            width: unset;
            height: unset;

            &__header {
                flex-direction: column;
                align-items: flex-start;
            }

            &__rule-area__rule {
                text-align: left;
            }
        }
    }
</style>
