import React, {Fragment, Component} from 'react';
import {connect} from "dva";
import { Link } from 'react-router-dom';
import { Modal, Form, Select, Button, Icon, Input, Col } from 'antd';
import styles from './index.module.css';
import lds from 'lodash';

@connect(({ loading, deploy, config }) => {
  return {
    gitBranchList: deploy.gitBranchList,
    gitCommitList: deploy.gitCommitList,
    gitTagList: deploy.gitTagList,
  }
})

class Ex1Info extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      appTdid: 0,
      loading: true,
      fetching: true,
      info: {},
      git_type: 'branch',
    }
  }

  componentDidMount() {
    // 清理 branchList
    this.props.dispatch({
      type: 'deploy/cleanBranchList',
    });

    this.setState({ 
      info: this.props.info,
    });

    if (this.props.info.GitType != undefined && this.props.info.GitType  != '') {
      this.setState({ 
        git_type: this.props.info.GitType,
      });
    }

    // 加载 branch 列表 和 branch[0] 的 commit 列表
    this.handleTempChange(this.props.info.Tid)

    // 加载 tag 列表
    this.handleTagList(this.props.info.Tid)
  }

  handleVersions = () => {
    this.props.dispatch({
      type: 'deploy/getAppVersion',
      payload: this.props.info.Tid,
    }).finally(() => this.setState({loading: false}));
  }

  // 返回所有的tag列表
  handleTagList = (tdid) => {
    this.props.dispatch({
      type: 'deploy/getGitTag',
      payload: this.props.info.Tid,
    }).finally(() => this.setState({loading: false}));
  }

  // 查branch内容
  handleTempChange = (tdid) => {
    if (tdid) {
      // 根据模板Tdid查branch
      this.setState({
        appTdid: tdid,
        loading: true,
      });
      const { dispatch } = this.props;
      dispatch({
        type: 'deploy/getGitBranch',
        payload: tdid,
      }).then(()=>{
        var tmpInfo = this.state.info;
        if (tmpInfo.TagBranch === undefined) {
          tmpInfo.TagBranch = this.props.gitBranchList[0]
        }
        if (tmpInfo.originTid != undefined && tmpInfo.originTid != tdid) {
          tmpInfo.TagBranch = this.props.gitBranchList[0]
        }
        this.setState({ 
          info: tmpInfo,
        });
        this.handleBrachChange(this.props.gitBranchList[0])
      });
    }
  }

  // 查commit 信息
  handleBrachChange = (values) => {
    if (values) {
      // 判断发布关联方式
      const items = {};
      items.aid  = this.state.appTdid;
      items.name = values
      // 根据项目id查branch 和 commit号
      this.setState({
        commitLoading: true,
      });
      const { dispatch } = this.props;
      dispatch({
        type: 'deploy/getGitCommit',
        payload: items,
      }).then(()=>{
        var tmpInfo = this.state.info;
        if (tmpInfo.Commit === undefined) {
          tmpInfo.Commit = this.props.gitCommitList[0]
        }

        if (tmpInfo.originTid !== undefined && tmpInfo.originTid != tmpInfo.Tid) {
          tmpInfo.Commit = this.props.gitCommitList[0]
        }
        this.setState({ 
          info: tmpInfo,
        });
      }).finally(() => this.setState({fetching: false}));
    }
  }

  fetchVersions = () => {
    if (this.state.git_type === 'branch') {
      this.setState({fetching: true});
      this.handleTempChange(this.props.info.Tid)
    }

    if (this.state.git_type === 'tag') {
      this.setState({loading: true});
      this.handleTagList(this.props.info.Tid)
    }
  };

  switchType = (v) => {
    var tmp = this.state.info;
    if (v==="branch") {
      tmp.TagBranch = this.props.gitBranchList[0]
    }
    if (v==="tag") {
      tmp.TagBranch = this.props.gitTagList[0]
    }

    this.setState({
      git_type: v, 
      info: tmp,
    });
  };

  switchExtra1 = (v) => {
    if (this.state.git_type === 'branch') {
      this.setState({fetching: true});
      this.handleBrachChange(v);
    }
    var tmp = this.state.info;
    tmp.TagBranch = v
    this.setState({
      info: tmp,
    })
  };

  handleSubmit = () => {
    this.setState({loading: true});
    const {git_type, info} = this.state;
    info.GitType = git_type;

    const { dispatch } = this.props;
    dispatch({
      type: 'deploy/deployAdd',
      payload: info,
    }).then(() => {
      this.setState({loading: false});
      this.props.onCancel();
    });
  };

  prePage = () => {
    this.props.prePage(this.state.info)
  };

  onDescInputChange = (e) => {
    var tmp = this.props.info;
    tmp['Desc'] = e.target.value;
    this.setState({ 
      info: tmp,
    });
  };

  handleCommitChange = (e) => {
    var tmp = this.props.info;
    tmp['Commit'] = e;
    this.setState({ 
      info: tmp,
    });
  };

  render() {
    const {info, fetching, git_type} = this.state;
    const {gitBranchList, gitCommitList, gitTagList} = this.props;

    return (
      <Form labelCol={{span: 6}} wrapperCol={{span: 14}}>
        <Form.Item required label="选择分支/标签/版本" help="根据网络情况，刷新可能会很慢，请耐心等待。">
          <Col span={19}>
            <Input.Group compact>
              <Select value={git_type} onChange={this.switchType} style={{width: 100}}>
                <Select.Option value="branch">Branch</Select.Option>
                <Select.Option value="tag">Tag</Select.Option>
              </Select>
              <Select
                showSearch
                style={{width: 320}}
                value={info.TagBranch}
                placeholder="请稍等"
                onChange={this.switchExtra1}
                filterOption={
                  (input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
                }
              >
                {git_type === 'branch' ? (
                  gitBranchList || {}).map(b => <Select.Option key={b} value={b}>{b}</Select.Option>
                ) : (
                  gitTagList || {}).map((item) => (
                    <Select.Option key={item} value={item}>{item}</Select.Option>
                  )
                )}
              </Select>
            </Input.Group>
          </Col>
          <Col span={4} offset={1} style={{textAlign: 'center'}}>
            {fetching ? <Icon type="loading" style={{fontSize: 18, color: '#1890ff'}}/> :
              <Button type="link" icon="sync" disabled={fetching} onClick={this.fetchVersions}>刷新</Button>
            }
          </Col>
        </Form.Item>
        {git_type === 'branch' && (
          <Form.Item required label="选择Commit">
            <Select value={info.Commit} 
              placeholder="请选择" onChange={this.handleCommitChange}>
              { gitBranchList ? gitCommitList.map(item => (
                <Select.Option key={item} value={item}>{item}</Select.Option>
              )) : null}
            </Select>
          </Form.Item>
        )}
        <Form.Item label="备注信息">
            <Input value={ info['Desc']} onChange={this.onDescInputChange} placeholder="请输入备注信息"/>
        </Form.Item>
        <Form.Item wrapperCol={{span: 14, offset: 6}}>
          <Button type="primary" disabled={this.state.fetching} onClick={this.handleSubmit}>提交</Button>
          <Button style={{marginLeft: 20}} onClick={this.prePage}>上一步</Button>
        </Form.Item>
      </Form> 
    )
  }
}

export default Ex1Info
