use anchor_lang::prelude::*;
use anchor_spl::token::{self, Token, TokenAccount, Transfer};

// 合约程序 ID（需替换为实际部署的 ID）
declare_id!("B7855fmCRGNGeFBh8R4tuDSdYkzeqpnqvxKcNMQsPmnE");

#[program]
pub mod sol_green {
    use super::*;

    // 初始化奖励配置
    pub fn initialize_reward_config(ctx: Context<InitializeRewardConfig>, base_reward: u64) -> Result<()> {
        let reward_config = &mut ctx.accounts.reward_config;
        reward_config.base_reward = base_reward;
        reward_config.admin = ctx.accounts.admin.key();
        msg!("奖励配置已初始化: base_reward={}, admin={}", base_reward, ctx.accounts.admin.key());
        Ok(())
    }

    // 发放环保奖励（仅管理员可调用）
    pub fn mint_green_reward(
        ctx: Context<MintGreenReward>,
        behavior_type: String,
        amount: u64,
    ) -> Result<()> {
        // 校验管理员权限
        require!(
            ctx.accounts.admin.key() == ctx.accounts.reward_config.admin,
            ErrorCode::Unauthorized
        );

        // 校验余额
        require!(
            ctx.accounts.admin_token_account.amount >= amount,
            ErrorCode::InsufficientBalance
        );

        // 代币转账逻辑
        let cpi_accounts = Transfer {
            from: ctx.accounts.admin_token_account.to_account_info(),
            to: ctx.accounts.user_token_account.to_account_info(),
            authority: ctx.accounts.admin.to_account_info(),
        };
        let cpi_program = ctx.accounts.token_program.to_account_info();
        let cpi_ctx = CpiContext::new(cpi_program, cpi_accounts);
        token::transfer(cpi_ctx, amount)?;

        // 记录奖励事件
        emit!(RewardMinted {
            user: ctx.accounts.user.key(),
            behavior_type: behavior_type.clone(),
            amount,
            timestamp: Clock::get()?.unix_timestamp,
        });

        msg!(
            "奖励已发放: user={}, behavior_type={}, amount={}",
            ctx.accounts.user.key(),
            behavior_type,
            amount
        );

        Ok(())
    }

    // 链上存证（记录行为哈希）
    pub fn record_proof(
        ctx: Context<RecordProof>,
        proof_hash: String,
        wallet_addr: String,
    ) -> Result<()> {
        let proof_record = &mut ctx.accounts.proof_record;
        proof_record.proof_hash = proof_hash;
        proof_record.wallet_addr = wallet_addr;
        proof_record.timestamp = Clock::get()?.unix_timestamp;
        proof_record.authority = ctx.accounts.authority.key();

        emit!(ProofRecorded {
            proof_hash: proof_record.proof_hash.clone(),
            wallet_addr: proof_record.wallet_addr.clone(),
            timestamp: proof_record.timestamp,
        });

        msg!("存证已记录: proof_hash={}, wallet={}", proof_record.proof_hash, proof_record.wallet_addr);

        Ok(())
    }
}

// 账户定义
#[derive(Accounts)]
pub struct InitializeRewardConfig<'info> {
    #[account(
        init,
        payer = admin,
        space = 8 + 8 + 32,
        seeds = [b"reward_config"],
        bump
    )]
    pub reward_config: Account<'info, RewardConfig>,
    #[account(mut)]
    pub admin: Signer<'info>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct MintGreenReward<'info> {
    #[account(
        seeds = [b"reward_config"],
        bump
    )]
    pub reward_config: Account<'info, RewardConfig>,
    #[account(mut)]
    pub admin: Signer<'info>,
    #[account(mut)]
    pub admin_token_account: Account<'info, TokenAccount>,
    #[account(mut)]
    pub user_token_account: Account<'info, TokenAccount>,
    pub user: AccountInfo<'info>,
    pub token_program: Program<'info, Token>,
}

#[derive(Accounts)]
pub struct RecordProof<'info> {
    #[account(
        init,
        payer = authority,
        space = 8 + 64 + 44 + 8 + 32,
        seeds = [b"proof", proof_hash.as_bytes()],
        bump
    )]
    pub proof_record: Account<'info, ProofRecord>,
    #[account(mut)]
    pub authority: Signer<'info>,
    pub system_program: Program<'info, System>,
}

// 数据结构
#[account]
pub struct RewardConfig {
    pub base_reward: u64,
    pub admin: Pubkey,
}

#[account]
pub struct ProofRecord {
    pub proof_hash: String,
    pub wallet_addr: String,
    pub timestamp: i64,
    pub authority: Pubkey,
}

// 事件
#[event]
pub struct RewardMinted {
    pub user: Pubkey,
    pub behavior_type: String,
    pub amount: u64,
    pub timestamp: i64,
}

#[event]
pub struct ProofRecorded {
    pub proof_hash: String,
    pub wallet_addr: String,
    pub timestamp: i64,
}

// 错误码
#[error_code]
pub enum ErrorCode {
    #[msg("Unauthorized access")]
    Unauthorized,
    #[msg("Insufficient balance")]
    InsufficientBalance,
    #[msg("Invalid proof hash")]
    InvalidProofHash,
}
