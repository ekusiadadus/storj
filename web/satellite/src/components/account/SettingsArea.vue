// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

<template>
    <div class="settings">
        <h1 class="settings__title" aria-roledescription="title">Account Settings</h1>
        <div class="settings__edit-profile">
            <div class="settings__edit-profile__row">
                <div class="settings__edit-profile__avatar">
                    <h1 class="settings__edit-profile__avatar__letter">{{ avatarLetter }}</h1>
                </div>
                <div class="settings__edit-profile__text">
                    <h2 class="profile-bold-text">Edit Profile</h2>
                    <h3 class="profile-regular-text">This information will be visible to all users</h3>
                </div>
            </div>
            <EditIcon
                class="edit-svg"
                @click="toggleEditProfileModal"
            />
        </div>
        <div class="settings__secondary-container">
            <div class="settings__secondary-container__change-password">
                <div class="settings__edit-profile__row">
                    <ChangePasswordIcon class="settings__secondary-container__img" />
                    <div class="settings__secondary-container__change-password__text-container">
                        <h2 class="profile-bold-text">Change Password</h2>
                        <h3 class="profile-regular-text">6 or more characters</h3>
                    </div>
                </div>
                <EditIcon
                    class="edit-svg"
                    @click="toggleChangePasswordModal"
                />
            </div>
            <div class="settings__secondary-container__email-container">
                <div class="settings__edit-profile__row">
                    <EmailIcon class="settings__secondary-container__img" />
                    <div class="settings__secondary-container__email-container__text-container">
                        <h2 class="profile-bold-text email">{{ user.email }}</h2>
                    </div>
                </div>
            </div>
        </div>
        <div class="settings__mfa">
            <h2 class="profile-bold-text">Two-Factor Authentication</h2>
            <p v-if="!user.isMFAEnabled" class="profile-regular-text">
                To increase your account security, we strongly recommend enabling 2FA on your account.
            </p>
            <p v-else class="profile-regular-text">
                2FA is enabled.
            </p>
            <div class="settings__mfa__buttons">
                <VButton
                    v-if="!user.isMFAEnabled"
                    label="Enable 2FA"
                    width="173px"
                    height="44px"
                    :on-press="enableMFA"
                    :is-disabled="isLoading"
                />
                <div v-else class="settings__mfa__buttons__row">
                    <VButton
                        class="margin-right"
                        label="Disable 2FA"
                        width="173px"
                        height="44px"
                        :on-press="toggleDisableMFAModal"
                        is-deletion="true"
                    />
                    <VButton
                        label="Regenerate Recovery Codes"
                        width="240px"
                        height="44px"
                        :on-press="generateNewMFARecoveryCodes"
                        is-blue-white="true"
                        :is-disabled="isLoading"
                    />
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

import { USER_ACTIONS } from '@/store/modules/users';
import { User } from '@/types/users';
import { APP_STATE_MUTATIONS } from "@/store/mutationConstants";

import VButton from '@/components/common/VButton.vue';

import ChangePasswordIcon from '@/../static/images/account/profile/changePassword.svg';
import EmailIcon from '@/../static/images/account/profile/email.svg';
import EditIcon from '@/../static/images/common/edit.svg';

// @vue/component
@Component({
    components: {
        EditIcon,
        ChangePasswordIcon,
        EmailIcon,
        VButton,
    },
})
export default class SettingsArea extends Vue {
    public isLoading = false;

    /**
     * Lifecycle hook after initial render where user info is fetching.
     */
    public mounted(): void {
        this.$store.dispatch(USER_ACTIONS.GET);
    }

    /**
     * Generates user's MFA secret and opens popup.
     */
    public async enableMFA(): Promise<void> {
        if (this.isLoading) return;

        this.isLoading = true;

        try {
            await this.$store.dispatch(USER_ACTIONS.GENERATE_USER_MFA_SECRET);
            this.toggleEnableMFAModal();
        } catch (error) {
            await this.$notify.error(error.message);
        }

        this.isLoading = false;
    }

    /**
     * Toggles generate new MFA recovery codes popup visibility.
     */
    public async generateNewMFARecoveryCodes(): Promise<void> {
        if (this.isLoading) return;

        this.isLoading = true;

        try {
            await this.$store.dispatch(USER_ACTIONS.GENERATE_USER_MFA_RECOVERY_CODES);
            this.toggleMFACodesModal();
        } catch (error) {
            await this.$notify.error(error.message);
        }

        this.isLoading = false;
    }

    /**
     * Toggles enable MFA modal visibility.
     */
    public toggleEnableMFAModal(): void {
        this.$store.commit(APP_STATE_MUTATIONS.TOGGLE_ENABLE_MFA_MODAL_SHOWN);
    }

    /**
     * Toggles disable MFA modal visibility.
     */
    public toggleDisableMFAModal(): void {
        this.$store.commit(APP_STATE_MUTATIONS.TOGGLE_DISABLE_MFA_MODAL_SHOWN);
    }

    /**
     * Toggles MFA recovery codes modal visibility.
     */
    public toggleMFACodesModal(): void {
        this.$store.commit(APP_STATE_MUTATIONS.TOGGLE_MFA_RECOVERY_MODAL_SHOWN);
    }

    /**
     * Opens change password popup.
     */
    public toggleChangePasswordModal(): void {
        this.$store.commit(APP_STATE_MUTATIONS.TOGGLE_CHANGE_PASSWORD_MODAL_SHOWN);
    }

    /**
     * Opens edit account info modal.
     */
    public toggleEditProfileModal(): void {
        this.$store.commit(APP_STATE_MUTATIONS.TOGGLE_EDIT_PROFILE_MODAL_SHOWN);
    }

    /**
     * Returns user info from store.
     */
    public get user(): User {
        return this.$store.getters.user;
    }

    /**
     * Returns first letter of user name.
     */
    public get avatarLetter(): string {
        return this.$store.getters.userName.slice(0, 1).toUpperCase();
    }
}
</script>

<style scoped lang="scss">
    .settings {
        position: relative;
        font-family: 'font_regular', sans-serif;
        padding-bottom: 70px;

        &__title {
            font-family: 'font_bold', sans-serif;
            font-size: 32px;
            line-height: 39px;
            color: #263549;
            margin: 40px 0 25px;
        }

        &__edit-profile {
            height: 66px;
            width: calc(100% - 80px);
            border-radius: 6px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 37px 40px;
            background-color: #fff;

            &__row {
                display: flex;
                justify-content: flex-start;
                align-items: center;
            }

            &__avatar {
                width: 60px;
                height: 60px;
                border-radius: 6px;
                display: flex;
                align-items: center;
                justify-content: center;
                background: #e8eaf2;
                margin-right: 32px;

                &__letter {
                    font-family: 'font_medium', sans-serif;
                    font-size: 16px;
                    line-height: 23px;
                    color: #354049;
                }
            }
        }

        &__secondary-container {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-top: 40px;

            &__change-password {
                height: 66px;
                border-radius: 6px;
                display: flex;
                justify-content: space-between;
                align-items: center;
                padding: 37px 40px;
                background-color: #fff;
                width: calc(48% - 80px);

                &__text-container {
                    margin-left: 32px;
                }
            }

            &__email-container {
                height: 66px;
                border-radius: 6px;
                display: flex;
                justify-content: flex-start;
                align-items: center;
                padding: 37px 40px;
                background-color: #fff;
                width: calc(48% - 80px);

                &__text-container {
                    margin-left: 32px;
                }
            }

            &__img {
                min-width: 60px;
                min-height: 60px;
            }
        }

        &__mfa {
            margin-top: 40px;
            padding: 40px;
            border-radius: 6px;
            background-color: #fff;

            &__buttons {
                margin-top: 20px;

                &__row {
                    display: flex;
                    align-items: center;
                }
            }
        }
    }

    .margin-right {
        margin-right: 15px;
    }

    .edit-svg {
        cursor: pointer;

        &:hover {

            .edit-svg__rect {
                fill: #2683ff;
            }

            .edit-svg__path {
                fill: white;
            }
        }
    }

    .input-container.full-input,
    .input-wrap.full-input {
        width: 100%;
    }

    .profile-bold-text {
        font-family: 'font_bold', sans-serif;
        color: #354049;
        font-size: 18px;
        line-height: 27px;
    }

    .profile-regular-text {
        margin: 10px 0;
        color: #afb7c1;
        font-size: 16px;
        line-height: 21px;
    }

    .email {
        word-break: break-all;
    }

    @media screen and (max-width: 1300px) {

        .profile-container {

            &__secondary-container {
                flex-direction: column;
                justify-content: center;

                &__change-password {
                    width: calc(100% - 80px);
                }

                &__email-container {
                    margin-top: 40px;
                    width: calc(100% - 80px);
                }
            }
        }
    }

    @media screen and (max-height: 825px) {

        .profile-container {
            height: 535px;
            overflow-y: scroll;

            &__secondary-container {
                margin-top: 20px;

                &__email-container {
                    margin-top: 20px;
                }
            }

            &__button-area {
                margin-top: 20px;
            }
        }
    }
</style>
